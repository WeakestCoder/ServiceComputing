# HW1 cli-basic
## 设计思路
使用flag包对命令行参数进行解析，用os，bufio.NewReader对文件、os.Stdin进行读取，输出使用os.Stdout.Write()，标准错误输出使用 fmt.Fprintf(os.Stderr, "error_message")。
  
  1.命令行检查 对结构体内的值进行判断，检查参数个数，参数值是否符合要求
  
  2.读写部分 通过参数判断是通过键盘读入还是文件读入，确定fin = bufio.NewReader(f)或者fin = bufio.NewReader(os.Stdin)。  
  读取时根据-f的存在与否确定分页方式，以'\f'作为分页方法，则使用input_byte, err := fin.ReadByte()进行读取，否则则使用crc, err := fin.ReadString('\n')进行读取
  
  3.管道部分  使用了 os/exec 包来建立用于进程间通信的管道。

## 测试部分
生成两个文件，l_file和f_file，l文件有1000行，每行为line n emmm， n代表行数'\n'用于换行，而f文件则是以'\f'作为换行 也是1000行  
1.input **./selpg -s 1 -e 1 l_file**
  
  output   
  ![image](https://raw.githubusercontent.com/WeakestCoder/ServiceComputing/master/HW1/screenshot/1.png)
  
2.input **./selpg -s 1 -e 1 l_file > output**  

  output 
  
3.input **./selpg -s 2 -e 1 l_file 2>error**

  output
4.input **./selpg -s 1 -e 2 l_file > output 2>error**  

  output
  
5.input **./selpg -s 1 -e 1 l_file >output 2>/dev/null**

  output文件输出跟4一样。没有错误消息打印到命令行 
  
6.input **./selpg -s 1 -e 1 -l 3 l_file**  

  output 

7.input **./selpg -s 1 -e 6 -f f_file**

  output
  
8.input **./selpg -s 1 -e 1 -d lp1 l_line**
  
  output：无打印机 无法测试
