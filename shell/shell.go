package shell

import (
	"github.com/driusan/kernel/asm"
	"github.com/driusan/kernel/filesystem"
	"github.com/driusan/kernel/process"
)

func Run() {
	proc := process.New()

	cons, err := proc.Open("/dev/cons")
	if err != nil {
		println(err.Error())
		return
	}
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
    pwd  - display current working directory
    exit - quit the shell
`))
				case "ns":
					Ns(cons, proc.Namespace)
				case "ls":
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
				case "exit":
					goto exit
				default:
					cons.Write([]byte("Unknown command: "))
					cons.Write([]byte(c))
					cons.Write([]byte(" (with arguments: \""))
					cons.Write([]byte(args))
					cons.Write([]byte("\")"))
				}
			}
			cmdSize = 0
			cmdEnd = 0
			cons.WriteRune('\n')
			cons.Write(prompt)
		case ' ':
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

func Ls(proc process.Process, cons filesystem.Writer, args string) error {
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
	if files != nil {
		for _, file := range files {
			cons.Write([]byte(file.Name()))
			cons.Write([]byte{' '})
		}
	}
	return nil
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
