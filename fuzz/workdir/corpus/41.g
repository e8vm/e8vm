struct a { a [4]byte; b byte }
func main() { 
	var as [4]a
	printUint(uint(&as[1])-uint(&as[0]))
}