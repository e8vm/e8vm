struct A { 
	a int;
	func p() { printInt(a) }
	func q(a int) { printInt(a) }
}
func main() { var a A; a.p(); a.a=33; a.p(); a.q(78) }