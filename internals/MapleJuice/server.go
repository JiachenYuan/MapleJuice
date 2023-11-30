package maplejuice

import (
	"bufio"
	"context"
	sdfs "cs425-mp/internals/SDFS"
	"cs425-mp/internals/global"
	pb "cs425-mp/protobuf"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

type MapleJuiceServer struct {
	pb.UnimplementedMapleJuiceServer
}

// var (
// 	MAPLE_INTERMEDIATE_FILES_FOLDER string
// 	JUICE_INTERMEDIATE_FILES_FOLDER string
// )

// func init() {
// 	usr, err := user.Current()
// 	if err != nil {
// 		fmt.Printf("Error getting user home directory: %v \n", err)
// 	}
// 	MAPLE_INTERMEDIATE_FILES_FOLDER = filepath.Join(usr.HomeDir, "MAPLE_INTERMEDIATE_FILES")

// 	err = os.MkdirAll(MAPLE_INTERMEDIATE_FILES_FOLDER, 0777)
// 	if err != nil {
// 		fmt.Println("Error creating Maple directory:", err)
// 		return
// 	}
// }

func (s *MapleJuiceServer) Maple(ctx context.Context, in *pb.MapleRequest) (*pb.MapleResponse, error) {
	if !sdfs.IsCurrentNodeLeader() {
		return nil, errors.New("not a leader, but received Maple command")
	}
	numMaples := in.NumMaples
	mapleExePath := in.MapleExePath
	prefix := in.SdfsIntermediateFilenamePrefix
	assignments := assignMapleWorkToWorkers(in.SdfsSrcDirectory, int(numMaples))
	err := sendMapleRequestToWorkers(assignments, mapleExePath, prefix)
	if err != nil {
		return nil, err
	}
	return &pb.MapleResponse{
		Success: true,
	}, nil
}

func (s *MapleJuiceServer) MapleExec(ctx context.Context, in *pb.MapleExecRequest) (*pb.MapleExecResponse, error) {
	files := in.Files
	mapleExePath := in.MapleExePath
	prefix := in.SdfsIntermediateFilenamePrefix

	sdfs.HandleGetFile(mapleExePath, mapleExePath)
	for _, file := range files {
		sdfs.HandleGetFile(file.Filename, file.Filename)
	}
	runExecutableFileOnInputFiles(mapleExePath, files, prefix)
	resp := &pb.MapleExecResponse{
		Success: true,
	}
	return resp, nil

}

func runExecutableFileOnInputFiles(mapleExePath string, fileLines []*pb.FileLines, prefix string) error {
	var wg sync.WaitGroup
	var mapleExeErrors []error
	var mut sync.Mutex
	for _, fileLine := range fileLines {
		wg.Add(1)
		go func(fileLine *pb.FileLines) {
			defer wg.Done()
			KVCollection, err := runExecutableFileOnSingleInputFile(mapleExePath, fileLine)
			if err != nil {
				mut.Lock()
				mapleExeErrors = append(mapleExeErrors, err)
				mut.Unlock()
			} else {
				err = appendAllIntermediateResultToSDFS(KVCollection, prefix)
				if err != nil {
					mut.Lock()
					mapleExeErrors = append(mapleExeErrors, err)
					mut.Unlock()
				}
			}
		}(fileLine)
	}
	wg.Wait()
	if len(mapleExeErrors) > 0 {
		return fmt.Errorf("some maple execution tasks failed: %v", mapleExeErrors)
	}

	fmt.Printf("Successfully finished excuting maple exe\n")
	return nil
}

func runExecutableFileOnSingleInputFile(mapleExePath string, fileLine *pb.FileLines) (map[string][]string, error) {
	KVCollection := make(map[string][]string)
	file := fileLine.Filename
	startLine := int(fileLine.Range.Start)
	endLine := int(fileLine.Range.End)
	fmt.Printf("Running maple executable on file %s, line %d to %d\n", file, startLine, endLine)
	currentLine := 0
	inputFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)
	lines := make([]string, 0)
	for scanner.Scan() {
		if currentLine >= startLine && currentLine <= endLine {
			line := file + "##" + scanner.Text()
			lines = append(lines, line)
			if currentLine > endLine {
				break
			}
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from input file:", err)
		return nil, err
	}
	cmd := exec.Command("python3", mapleExePath)

	// Map exe's input will have information about which sdfs file this line is coming from before the ## sign
	cmd.Stdin = strings.NewReader(strings.Join(lines, "\n"))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing script on line %d: %s\n", currentLine, err)
		return nil, err
	}
	kvPairs := strings.Split(string(output), "\n")
	for _, kvPair := range kvPairs {
		kv := strings.Split(kvPair, ":")
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		value := kv[1]
		KVCollection[key] = append(KVCollection[key], value)
	}

	return KVCollection, nil
}

func (s *MapleJuiceServer) JuiceExec(ctx context.Context, in *pb.JuiceExecRequest) (*pb.JuiceExecResponse, error) {
	// Extract request fields
	juiceProgram := in.JuiceProgram
	dstFileName := in.DstFilename

	// Get file from SDFS
	sdfs.HandleGetFile(juiceProgram, juiceProgram)
	for _, file := range in.InputIntermFiles {
		sdfs.HandleGetFile(file, file)
	}

	// Create a temp file holding local aggregate results for all assigned keys

	f, err := os.CreateTemp("", "juice_local_result")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	// todo: make the parsing job concurrent, the file IO can be sequential and that's fine
	for _, inputFilename := range in.InputIntermFiles {
		file, err := os.Open(inputFilename)
		if err != nil {
			fmt.Printf("unable to open intermediate file input %s: %v\n", inputFilename, err)
			return nil, err
		}

		key := ""
		values := "" // todo: value set might be too big, move it to disk if possible
		// Read file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// Split the line into key and value
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				key = parts[0]
				value := parts[1]

				values += value + "::"
			} else {
				fmt.Println("Invalid line format:", line)
				return nil, errors.New("Invalid line format:" + line)
			}
		}

		valuesStr := values[:len(values)-2] // remove the last deliemeter (::)
		programInputStr := fmt.Sprintf("%s:%s", key, valuesStr)
		// Give value set to the juice task executable
		cmd := exec.Command("python3", juiceProgram)
		cmd.Stdin = strings.NewReader(programInputStr)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error executing script on line %s: %s\n", programInputStr, err)
			return nil, err
		}

		// Write the parsed key: [values set] into the temp file
		f.Write(output)
	}

	// Append (create if necessary) temp file content to destination global file
	data, err := os.ReadFile(f.Name())
	if err != nil {
		fmt.Printf("cannot read the temporary file: %v\n", err)
	}
	sdfs.HandleAppendFile(dstFileName, string(data))

	return &pb.JuiceExecResponse{
		Success: true,
	}, nil
}

// Only leader should process at this endpoint
func (s *MapleJuiceServer) Juice(ctx context.Context, in *pb.JuiceRequest) (*pb.JuiceResponse, error) {
	if !sdfs.IsCurrentNodeLeader() {
		return nil, errors.New("not a leader, but received Juice command")
	}

	// Extract request fields
	juiceProgram := in.JuiceExecName
	numJuicer := int(in.NumJuicer)
	filePrefix := in.Prefix
	dstFileName := in.DestName
	deleteInputAfter := in.DeleteInput
	useRangePartition := in.IsRangePartition

	// var vmToInputFiles map[string]map[string]global.Empty
	vmToInputFiles := createKeyAssignmentForJuicers(numJuicer, filePrefix, useRangePartition)
	err := dispatchJuiceTasksToVMs(vmToInputFiles, juiceProgram, dstFileName)

	if err != nil {
		return nil, err
	}

	if deleteInputAfter {
		// Delete all the input files
		for _, fileSet := range vmToInputFiles {
			for sdfsFilename := range fileSet {
				global.MemTable.DeleteFile(sdfsFilename)
			}
		}
	}

	return &pb.JuiceResponse{
		Success: true,
	}, nil
}

func appendAllIntermediateResultToSDFS(KVCollection map[string][]string, prefix string) error {
	// Iterate over the directory entries and delete each file.
	for key, values := range KVCollection {
		var content string
		for _, v := range values {
			content += fmt.Sprintf("%s:%s\n", key, v)
		}
		sdfsIntermediateFileName := fmt.Sprintf("%s_%s", prefix, key)
		fmt.Printf("Trying to appended to SDFS file %s\n", sdfsIntermediateFileName)
		sdfs.HandleAppendFile(sdfsIntermediateFileName, content)
	}

	return nil
}

func StartMapleJuiceServer() {
	// start listening for incoming connections
	lis, err := net.Listen("tcp", ":"+global.MAPLE_JUICE_PORT)
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterMapleJuiceServer(s, &MapleJuiceServer{})
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v\n", err)
	}
}
