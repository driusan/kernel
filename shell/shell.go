package shell

import (
	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/executable"
	"github.com/driusan/kernel/filesystem"
	"github.com/driusan/kernel/memory"
	"github.com/driusan/kernel/process"
	"github.com/driusan/kernel/terminal"
)

func Run() {
	proc := process.New()

	cons, err := proc.Open("/dev/cons")
	if err != nil {
		println(err.Error())
		return
	}

	// STDIN, STDOUT, and STDERR all default to cons.
	proc.FDs = []filesystem.File{cons, cons, cons}
	prompt := []byte{'>', ' '}

	cons.Write([]byte(`
Welcome to the shell.
This is not yet implemented very well, but feel free to poke around.

Type help for a list of valid commands.
`))
	cons.Write(prompt)

	// This isn't working. There's a bug somewhere in libg causes it to
	// result in an index out of range exception, so for now just use
	// an array.
	//cmd := make([]byte)
	var cmd [1024]byte
	cmdSize := 0
	cmdEnd := 0
	for {
		asm.HLT()
		chr, err := cons.ReadByte()
		if err != nil {
			println(err.Error())
			break
		}

		switch chr {
		case '\n':
			if cmdSize > 0 {
				var args string
				if cmdEnd == 0 {
					cmdEnd = cmdSize
					args = ""
				} else {
					args = string(cmd[cmdEnd+1 : cmdSize])
				}
				c := string(cmd[0:cmdEnd])
				cons.WriteRune('\n')

				switch c {
				case "help":
					cons.Write([]byte(`
Valid commands:
    help - display this help
    ns   - display process namespace
    ls   - list files
    cd   - change current working directory
    cat file - print file to the screen
    mem  - display memory usage
    pwd  - display current working directory
    exit - quit the shell
`))
				case "ns":
					// This should be a separate command once forking
					// is set up
					Ns(cons, proc.Namespace)
				case "ls":
					// This should be a separate command once processes
					// are setup
					err := Ls(proc, cons, args)
					if err != nil {
						cons.Write([]byte(err.Error()))

					}
				case "cd":
					err := proc.Cwd(filesystem.Path(args))
					if err != nil {
						cons.Write([]byte(err.Error()))
					}
				case "pwd":
					if proc.Wd != "" {
						cons.Write([]byte(proc.Wd))
					} else {
						cons.Write([]byte("No current working directory"))
					}
				case "mem":
					// This should be part of /dev, not a command
					alloc, free, err := memory.MemStats()
					if err != nil {
						cons.Write([]byte(err.Error()))
					} else {
						cons.Write([]byte("Allocated pages: "))
						terminal.PrintDec(alloc)
						cons.Write([]byte(" Free pages: "))
						terminal.PrintDec(free)
					}
				case "cat":
					// This should be a separate command once userspace
					// processes are setup
					infile, err := proc.Open(filesystem.Path(args))
					if err != nil {
						cons.Write([]byte(err.Error()))
					} else {
						err = Cat(infile, cons)
						if err != nil {
							cons.Write([]byte(err.Error()))
						}
					}
				case "exit":
					goto exit
				default:
					// The indentation here is weird, the flow
					// should be cleaned up a little. I don't like
					// else statements.
					file, err := proc.Open(filesystem.Path(cmd[:cmdEnd]))
					if err != nil {
						cons.Write([]byte("Unknown command: "))
						cons.Write([]byte(c))
						cons.Write([]byte(" (with arguments: \""))
						cons.Write([]byte(args))
						cons.Write([]byte("\")"))
					} else {
						err = executable.Run(file, &proc)
						if err != nil {
							cons.Write([]byte(err.Error()))
						}
						// this should be done with a defer, once enough
						// of runtime is implemented to use defer...
						file.Close()
					}
				}
			}
			cmdSize = 0
			cmdEnd = 0
			cons.WriteRune('\n')
			cons.Write(prompt)
		case ' ':
			// mark the border between a command and its arguments
			// if applicable, but otherwise treat it the same as any
			// other character
			if cmdEnd == 0 {
				cmdEnd = cmdSize
			}
			fallthrough
		default:
			cmd[cmdSize] = chr
			cmdSize++
			cons.WriteRune(rune(chr))
		}
	}
exit:
	println("Leaving the shell.")
}

func Ls(proc process.Process, cons filesystem.File, args string) error {
	var dirName filesystem.Path
	if len(args) == 0 {
		dirName = proc.Wd
	} else {
		if args[0] == '/' {
			dirName = filesystem.Path(args)
		} else {
			dirName = proc.Wd + "/" + filesystem.Path(args)
		}
	}
	if dirName == "" {
		return filesystem.FilesystemError("No current working directory")
	}

	d, err := proc.Open(dirName)
	if err != nil {
		return err
	}

	dir, err := d.AsDirectory()
	if err != nil {
		return err
	}

	files := dir.Files()

	for name, file := range files {
		cons.Write([]byte(name))
		if file.IsDirectory() {
			cons.WriteRune('/')
		}
		cons.Write([]byte{' '})
	}
	return nil
}

func Cat(stdin, stdout filesystem.File) error {
	output := make([]byte, 4096)
	for {
		n, err := stdin.Read(output)
		/*
			need to implement __go_interface_value_compare for this switch to work
			as it should. For now just use a hack of checking the Error() string
			value after making sure err isn't nil
		*/
		if err == nil {
			stdout.Write(output[0:n])
		} else {
			switch err.Error() { // err
			case "End of file": //filesystem.EOF:
				stdout.Write(output[0:n])
				return nil
			case "": // nil
				stdout.Write(output[0:n])
			default:
				return err
			}
		}
	}
}

func Ns(cons filesystem.File, ns process.Namespace) {
	for path, handler := range ns {
		if path != "" {
			cons.Write([]byte(path))
			cons.WriteRune('\t')
			cons.Write([]byte(handler.Type()))
			cons.WriteRune('\n')
		}
	}
}
