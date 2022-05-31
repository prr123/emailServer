package iploc

import (
	"fmt"
	"os"
	"strconv"
	"io"
	"unsafe"
)

type Countries struct {
	Short string
	Name string
	Num int
}

type Iplocrec struct {
	Ipst uint32
	Ipend uint32
	CC []byte
	CCstr string
	Name string
	Num int
}

type rdstate struct {
	rdbuf [Blk]byte
	rdidx int
	nb int
	ibnxt int
	ibrec int
	iof int64
	niof int64
	inxt int
	state int
	Totlin int
	Maxblk int
}

const Blk = 4096
const Srecsiz = 8
const Ipd = 256
const Slblk = 66

func Str2btoint(numstr string)(inum uint32) {
	var nb [2]byte
	var i16num uint16
	if len(numstr) < 2 {
		inum = 0
		return inum
	}
	if len(numstr) >2 {
		inum = 1
		return inum
	}

	xs := []byte(numstr)
	for i:=0; i< 2; i++ {
        nb[i] = xs[1-i]
    }
    i16num = *(*uint16)(unsafe.Pointer(&nb[0]))
	inum = uint32(i16num)
	return inum
}

func Intto2bstr(ccinp uint32)(ccstr string){

    bufptr := (*[2]byte)(unsafe.Pointer(&ccinp))
    xb:= (*bufptr)[0:2]
    ibyte := xb[1]
    xb[1] = xb[0]
    xb[0] = ibyte
    ccstr = string(xb)
    return ccstr
}

func B2tostr(bin []byte)(ccstr string){

    bufptr := (*[]byte)(unsafe.Pointer(&bin[0]))
    xb:= (*bufptr)[0:2]
    ibyte := xb[1]
    xb[1] = xb[0]
    xb[0] = ibyte
    ccstr = string(xb)
    return ccstr
}



func Int2tobyte (xin [2]uint32)(bs []byte) {
	bufptr := (*[8]byte)(unsafe.Pointer(&xin[0]))
	bs = (*bufptr)[:]
	return bs
}

func Bytetoint2 (bs []byte) (zout *[2]uint32) {
	if len(bs) !=8 {
		return nil
	}
	zout = (*[2]uint32)(unsafe.Pointer(&bs[0]))
	return zout
}


func Inttoipstr (xin uint32) (ipstr string) {
	var bs [4]byte
    bs = *(*[4]byte)(unsafe.Pointer(&xin))

	ipstr = fmt.Sprintf("%d.%d.%d.%d",bs[3],bs[2],bs[1],bs[0])

    return ipstr
}

func Byte4toint_end (bs []byte) (zout *uint32) {
    var tbyt [4]byte
    if len(bs) !=4 {
        return nil
    }
    tbyt[0] = bs[3]
    tbyt[1] = bs[2]
    tbyt[2] = bs[1]
    tbyt[3] = bs[0]

    zout = (*uint32)(unsafe.Pointer(&tbyt[0]))
    return zout
}

func Intto4byte_end (xin uint32)(bs []byte) {
    var tbyt [4]byte
    bu := (*[4]byte)(unsafe.Pointer(&xin))

    tbyt[0] = bu[3]
    tbyt[1] = bu[2]
    tbyt[2] = bu[1]
    tbyt[3] = bu[0]

    bs = tbyt[:]
    return bs

}
func Byte4toint (bs []byte) (zout uint32) {
    zout = *(*uint32)(unsafe.Pointer(&bs[0]))
    return zout
}

func Intto4byte (xin uint32)(bs []byte) {
    bu := (*[4]byte)(unsafe.Pointer(&xin))
    bs = bu[:]
    return bs
}

func Cvt_int(xin uint32) (idx uint32) {
    var tbyt [4]byte
    bu := (*[4]byte)(unsafe.Pointer(&xin))

    tbyt[0] = bu[3]
    tbyt[1] = bu[2]
    tbyt[2] = bu[1]
    tbyt[3] = bu[0]

    zout := (*uint32)(unsafe.Pointer(&tbyt[0]))

    idx = *zout
    return idx
}

func Cvt2_int(xin uint32) (idx uint32) {
    bu := (*[4]byte)(unsafe.Pointer(&xin))
    tbyt := (*[4]byte)(unsafe.Pointer(&idx))

    tbyt[0] = bu[3]
    tbyt[1] = bu[2]
    tbyt[2] = bu[1]
    tbyt[3] = bu[0]

    return idx
}

func Cvt_from_sbrec(inbuf [8]byte) (ip uint32, cstr string){
	iptr:= (*[4]byte)(unsafe.Pointer(&ip))
	for i:=0; i<4; i++ {
		iptr[i] = inbuf[i]
	}
	sptr:= (*[2]byte)(unsafe.Pointer(&inbuf[4]))
	xb := (*sptr)[0:2]
	xbp := xb[:]
	cstr = string(xbp)
	return ip, cstr
}

func Cvt_to_sbrec(iprec Iplocrec)(wb [8]byte) {
	var wptr *[4]byte
	var sptr []byte
	wptr = (*[4]byte)(unsafe.Pointer(&iprec.Ipst))

	wb[0] = wptr[0]
	wb[1] = wptr[1]
	wb[2] = wptr[2]
	wb[3] = wptr[3]
	sptr = []byte(iprec.CCstr)
	wb[4] = sptr[0]
	if len(iprec.CCstr) > 1 {wb[5] = sptr[1] }
	return wb
}

func Rdiptxtfile(fil *os.File, irec int)(iprec Iplocrec, err error) {

	if fil == nil {
		err = fmt.Errorf("error did not supply a fd!")
		return iprec, err
	}




	return iprec, nil
}


func Newrdstate()(rd rdstate){

	rd.rdidx = 0
	rd.nb = 0
	rd.ibnxt =0
	rd.ibrec = 0
	rd.iof = 0
	rd.inxt =0
	rd.state = 0
	rd.Totlin =0
	return rd
}

func (rd *rdstate)Rdipnext(fil *os.File)(iprec Iplocrec, err error) {

	if fil == nil {
		err = fmt.Errorf("error did not supply a fd!")
		return iprec, err
	}

//	buf := rd.rdbuf[:]

//	istate:=rd.state

//	fmt.Println("blocks: ", rd.ibnxt, rd.Maxblk, rd.inxt)
	for ib:=rd.ibnxt; ib<rd.Maxblk; ib++ {
//	fmt.Println("start block: ", ib)
		ip1str :=""
		ip2str := ""
		ccstr:= ""

		ist:= rd.inxt
		buf := rd.rdbuf[:]
//		fmt.Println("next buf: ", ib, string(buf[ist:ist+40]))

		wdst := rd.inxt
		istate := 0

		for i:=rd.inxt; i<rd.nb; i++ {
//		for i:=rd.inxt; i<120; i++ {
			b:= rd.rdbuf[i]

			switch b {
			case '"':
				switch istate {
				case 0:
					istate = 1
					wdst = i+1

				case 1:
					istate =2
					ip1str = string(buf[wdst:i])
					ipt, err :=strconv.Atoi(ip1str)

					if err != nil {
						fmt.Println("err atoi line: ", rd.Totlin, err)
					}
					iprec.Ipst = uint32(ipt)
//					fmt.Println("ip1: ", iprec.Ipst, "| ", ip1str)

				case 3:
					istate = 4
					wdst = i+1

				case 4:
					istate =5
					ip2str = string(buf[wdst:i])
					ipt, err :=strconv.Atoi(ip2str)
					if err != nil {
						fmt.Println("err atoi line: ", rd.Totlin, err)
					}
					iprec.Ipend = uint32(ipt)
//					fmt.Println("ipe: ", iprec.Ipend, "| ", ip2str)

				case 6:
					istate = 7
					wdst = i+1

				case 7:
					istate =8
					ccstr = string(buf[wdst:i])
					if i-wdst > 2 {
						err = fmt.Errorf("error line: ", rd.Totlin, "cc not 2 char")
						return iprec, err
 					}
					copy (iprec.CC[:], buf[wdst:i])
					iprec.CCstr = ccstr
//					fmt.Println("cc: ", iprec.CC, "| ", iprec.CCstr)

				case 9:
					istate = 10
					wdst = i+1

				case 10:
					istate = 11
					iprec.Name = string(buf[wdst:i])
					iprec.Num = rd.Totlin
					rd.niof = rd.iof + int64(i+1)

					rd.inxt = i+1
					rd.Totlin++
					rd.state = istate

//					fmt.Println(" line f[", rd.Totlin, "] ", string(buf[ist:i+1]))

					return iprec, nil


				default:
					fmt.Println("error parsing line: ",rd.Totlin, " pos: ", i ," state: ", istate)
					fmt.Println(" line d[", rd.Totlin, "] ", ist, i, string(buf[ist:i]))
				}

			case ',' :
				switch istate {
					case 2,5,8:
						istate++

					case 10:

					default:
						err = fmt.Errorf("error parsing line: ", rd.Totlin, " pos: ", i, " state: ", istate)
						fmt.Println(" line com[", rd.Totlin, "] ", string(buf[ist:i]))
				}

			case '\n':
				switch istate {
				case 0:
					ist = i+1

				case 11:
					ist = i+1
					istate =0

				default:
					fmt.Println("error cr parsing line: ",rd.Totlin, " pos: ", i, " state: ", istate)
					fmt.Println(" line [", rd.Totlin, "] ", string(buf[ist:i]))
				}

			} // switch

		} //i


//		fmt.Println("iof: ", rd.iof, rd.niof)
		if (rd.iof == rd.niof)&&(rd.iof>0) {
			return iprec, io.EOF
		}
		rd.iof = rd.niof

		rd.nb, err = fil.ReadAt(buf, rd.iof)
//		fmt.Println("read bytes: ", rd.nb, err)
/*
		if (err == io.EOF)&&(rd.nb<1) {
			return iprec, err
		}
*/
		if (err != nil) && (err != io.EOF) {
			err = fmt.Errorf(" error reading block: ", ib, " err: ", err)
			return iprec, err
		}
		rd.inxt = 0

//		fmt.Println("read block: ", ib, " nb: ", rd.nb, " iof: ", rd.iof)
//		fmt.Println("block buf: ", string(buf[ist:ist+40]))

	} //ib

	err = fmt.Errorf("error:  did not finish reading csv file!")
	return iprec, err
}


func Rec_prnt (iprec Iplocrec) {
//	ccstring := B2tostr(iprec.CC)
	fmt.Println("rec num: ", iprec.Num, " ip st: ", iprec.Ipst, " ip end: ", iprec.Ipend, " cc: ", iprec.CCstr," Name: ", iprec.Name)
}

