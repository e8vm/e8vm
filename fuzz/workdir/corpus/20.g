struct A { b byte; a int }
func main() { var a A; var pa=&a; 
	pa.a = 33; pa.b = byte(7);
	printInt(pa.a); printInt(int(pa.b))
}