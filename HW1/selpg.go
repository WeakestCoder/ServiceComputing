package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

var progname string = ""

type selpg struct {
	s           int
	e           int
	in_filename string
	page_len    int
	page_type   bool
	print_dest  string
}

func process_input(temp selpg) {
	var fin_ptr *os.File
	var fin *bufio.Reader
	var fout *bufio.Writer
	var stdinpipe io.WriteCloser
	var cmd *exec.Cmd
	var line_ctr int
	var page_ctr int

	if temp.in_filename != "" {
		f, err := os.Open(temp.in_filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s:open input file fail!\n", temp.in_filename)
		}
		fin_ptr = f
		fin = bufio.NewReader(f)
	} else {
		fin = bufio.NewReader(os.Stdin)
	}

	if temp.print_dest != "" {
		/*to be completed*/
		var dest_flag string = fmt.Sprintf("-d%s", temp.print_dest)
		cmd = exec.Command("lp", dest_flag)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			if fin_ptr != nil {
				fin_ptr.Close()
			}
			fmt.Fprintf(os.Stderr, "%s: could not open pipe to \"%s\"\n", progname, fmt.Sprintf("lp %s", dest_flag))
		}

		stdinpipe = stdin
		fout = bufio.NewWriter(stdin)
		err = cmd.Start()

		if err != nil {
			if fin_ptr != nil {
				fin_ptr.Close()
			}
			stdin.Close()
			fmt.Fprintf(os.Stderr, "%s: could not start a printer when run \"%s\"\n", progname, fmt.Sprintf("lp %s", dest_flag))
		}
	} else {
		fout = bufio.NewWriter(os.Stdout)
	}

	if !temp.page_type {
		line_ctr = 0
		page_ctr = 1
		for true {
			crc, err := fin.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			line_ctr++
			if line_ctr > temp.page_len {
				page_ctr++
				line_ctr = 1
			}

			if page_ctr >= temp.s && page_ctr <= temp.e {
				_, err := fout.Write([]byte(crc))
				if err != nil {
					panic(err)
				}
				fout.Flush()
			}
		}
	} else {
		page_ctr = 1
		for true {
			input_byte, err := fin.ReadByte()
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			if input_byte == '\f' {
				page_ctr++
			}
			if page_ctr >= temp.s && page_ctr <= temp.e {
				err := fout.WriteByte(input_byte)
				if err != nil {
					if fin_ptr != nil {
						fin_ptr.Close()
					}
					if temp.print_dest != "" {
						stdinpipe.Close()
					}
					panic(err)
				}
				fout.Flush()
			}
		}
	}

	if page_ctr < temp.s {
		fmt.Fprintf(os.Stderr, "%s: start_page (%d) greater than total pages (%d),no output written\n", progname, temp.s, page_ctr)
	} else if page_ctr < temp.e {
		fmt.Fprintf(os.Stderr, "%s: end_page (%d) greater than total pages (%d),less output than expected\n", progname, temp.s, page_ctr)
	}

	fout.Flush()

	if temp.print_dest != "" {
		stdinpipe.Close()
		err := cmd.Wait()
		if err != nil {
			if fin_ptr != nil {
				fin_ptr.Close()
			}
			fmt.Fprintf(os.Stderr, "%s: complete printing error with the printer \n", progname)
		}
	}
	if fin_ptr != nil {
		fin_ptr.Close()
	}
	fmt.Fprintf(os.Stderr, "%s: done\n", progname)
}

func usage() {
	fmt.Fprintf(os.Stderr, "\nUSAGE:selge -s start_page -e end_page [ -f | -l lines_per_page ] [ -ddest ] [ in_filename ]\n")
}

func main() {
	progname = os.Args[0]
	var temp selpg
	flag.IntVar(&temp.s, "s", -1, "specify start_page to use.(>=1).")
	flag.IntVar(&temp.e, "e", -1, "specify end_page to use.(>=s).")
	flag.IntVar(&temp.page_len, "l", 72, "specify the length of page.")
	flag.BoolVar(&temp.page_type, "f", false, "specify the type of page end")
	flag.StringVar(&temp.print_dest, "d", "", "specify print destination")
	flag.Usage = usage
	flag.Parse()
	/*fmt.Printf("temp.s = %d\n", temp.s)
	fmt.Printf("temp.e = %d\n", temp.e)
	fmt.Printf("temp.page_len = %d\n", temp.page_len)
	fmt.Printf("dest = %s\n", temp.print_dest)
	fmt.Printf("temp.type = %t\n", temp.page_type)*/
	if temp.s == -1 || temp.e == -1 || temp.s > temp.e || temp.s < 1 || temp.e < 1 || (temp.page_len != 72 && temp.page_type == true) || len(flag.Args()) > 1 {
		flag.Usage()
		return
	}
	if len(flag.Args()) == 1 {
		temp.in_filename = flag.Args()[0]
	}

	f, _ := os.Open(temp.in_filename)

	if f == nil {
		fmt.Fprintf(os.Stderr, "file %s doesn't exist\n", temp.in_filename)
		flag.Usage()
		return
	}

	process_input(temp)
}
