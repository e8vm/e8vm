package arch

import (
	"testing"

	"math/rand"
)

func TestInstImm(t *testing.T) {
	m := newPhyMemory(PageSize * 32)
	cpu := newCPU(m, nil, new(instImm), 0)

	tst := func(op, s, v, im, d, res uint32) {
		cpu.Reset()

		in := (op & 0xff) << 24
		in |= (s & 0x7) << 18
		in |= (d & 0x7) << 21
		in |= im & 0xffff

		m.WriteU32(InitPC, in)

		cpu.regs[s] = v
		e := cpu.Tick()
		if e != nil {
			t.Fatal("unexpected exception")
		}

		got := cpu.regs[d]
		if got != res {
			t.Fatalf("got 0x%08x, expect 0x%08x", got, res)
		}
	}

	twr := func(op, s, v, im, d, w uint32) {
		cpu.Reset()

		in := (op & 0xff) << 24
		in |= (s & 0x7) << 18
		in |= (d & 0x7) << 21
		in |= im & 0xffff

		m.WriteU32(InitPC, in)
		cpu.regs[s] = v
		cpu.regs[d] = w
		e := cpu.Tick()
		if e != nil {
			t.Fatal("unexpected exception")
		}
	}

	tf := func(op uint32, f func(v, im uint32) uint32) {
		for i := 0; i < 100; i++ {
			s := uint32(rand.Intn(5))
			d := uint32(rand.Intn(5))
			v := uint32(rand.Int63())
			im := uint32(rand.Int63()) & 0xffff

			exp := f(v, im)

			tst(op, s, v, im, d, exp)
		}
	}

	sext := func(i uint32) uint32 {
		return uint32((int32(i << 16)) >> 16)
	}

	tf(ADDI, func(v, im uint32) uint32 { return v + sext(im) })
	tf(SLTI, func(v, im uint32) uint32 {
		se := sext(im)
		if int32(v) < int32(se) {
			return 1
		}
		return 0
	})
	tf(ANDI, func(v, im uint32) uint32 { return v & im })
	tf(ORI, func(v, im uint32) uint32 { return v | im })
	tf(XORI, func(v, im uint32) uint32 { return v ^ im })
	tf(ADDUI, func(v, im uint32) uint32 { return v + (im << 16) })

	for i := 0; i < 100; i++ {
		addr := uint32(PageSize * 10)
		addr += uint32(rand.Int63()) % PageSize * 4
		offset := uint32(rand.Int63()) % PageSize
		offset -= offset % 4
		offset -= PageSize / 2
		w := uint32(rand.Int63())

		e := m.WriteU32(addr+offset, w)
		if e != nil {
			t.Fatal("write fail")
		}
		s := uint32(rand.Intn(5))
		d := uint32(rand.Intn(5))

		tst(LW, s, addr, offset, d, w)
	}

	for i := 0; i < 100; i++ {
		addr := uint32(PageSize * 10)
		addr += uint32(rand.Int63()) % PageSize * 4
		offset := uint32(rand.Int63()) % PageSize
		offset -= PageSize / 2
		b := byte(rand.Int63())

		e := m.WriteU8(addr+offset, b)
		if e != nil {
			t.Fatal("write fail")
		}
		s := uint32(rand.Intn(5))
		d := uint32(rand.Intn(5))

		w := uint32(int32(uint32(b)<<24) >> 24)
		tst(LB, s, addr, offset, d, w)
		tst(LBU, s, addr, offset, d, uint32(b))
	}

	for i := 0; i < 100; i++ {
		addr := uint32(PageSize * 10)
		addr += uint32(rand.Int63()) % PageSize * 4
		offset := uint32(rand.Int63()) % PageSize
		offset -= offset % 4
		offset -= PageSize / 2
		w := uint32(rand.Int63())

		s := uint32(rand.Intn(5))
		d := uint32(rand.Intn(5))
		for d == s {
			d = uint32(rand.Intn(5))
		}

		twr(SW, s, addr, offset, d, w)

		got, e := m.ReadU32(addr + offset)
		if e != nil {
			t.Fatal("read fail")
		}

		if got != w {
			t.Fatalf("expect 0x%08x, got 0x%08x", w, got)
		}
	}

	for i := 0; i < 100; i++ {
		addr := uint32(PageSize * 10)
		addr += uint32(rand.Int63()) % PageSize * 4
		offset := uint32(rand.Int63()) % PageSize
		offset -= PageSize / 2
		b := byte(rand.Int63())

		s := uint32(rand.Intn(5))
		d := uint32(rand.Intn(5))
		for d == s {
			d = uint32(rand.Intn(5))
		}

		twr(SB, s, addr, offset, d, uint32(b))

		got, e := m.ReadU8(addr + offset)
		if e != nil {
			t.Fatal("read fail")
		}

		if got != b {
			t.Fatalf("expect 0x%08x, got 0x%08x", b, got)
		}
	}
}
