package tpcc

import (
	"context"
	"fmt"
	"time"

	"github.com/siddontang/go-tpc/pkg/load"
)

const (
	maxItems              = 100000
	stockPerWarehouse     = 100000
	districtPerWarehouse  = 10
	customerPerWarehouse  = 30000
	customerPerDistrict   = 3000
	orderPerWarehouse     = 30000
	historyPerWarehouse   = 30000
	newOrderPerWarehouse  = 9000
	orderLinePerWarehouse = 300000
	minOrderLinePerOrder  = 5
	maxOrderLinePerOrder  = 15
)

func (w *Workloader) loadItem(ctx context.Context, tableID int) error {
	s := w.base.GetState(ctx)
	hint := fmt.Sprintf("INSERT INTO item%d (i_id, i_im_id, i_name, i_price, i_data) VALUES ", tableID)
	l := load.NewBatchLoader(s.Conn, hint)

	for i := 0; i < maxItems; i++ {
		s.Buf.Reset()
		iImID := randInt(s.R, 1, 10000)
		iPrice := float64(randInt(s.R, 100, 10000)) / float64(100.0)
		iName := randChars(s.R, s.Buf, 14, 24)
		iData := randOriginalString(s.R, s.Buf)

		v := fmt.Sprintf(`(%d, %d, '%s', %f, '%s')`, i+1, iImID, iName, iPrice, iData)

		if err := l.InsertValue(ctx, v); err != nil {
			return err
		}
	}

	return l.Flush(ctx)
}

func (w *Workloader) loadWarhouse(ctx context.Context, tableID int, warehouse int) error {
	s := w.base.GetState(ctx)
	hint := fmt.Sprintf("INSERT INTO warehouse%d (w_id, w_name, w_street_1, w_street_2, w_city, w_state, w_zip, w_tax, w_ytd) VALUES ", tableID)
	l := load.NewBatchLoader(s.Conn, hint)

	wName := randChars(s.R, s.Buf, 6, 10)
	wStree1 := randChars(s.R, s.Buf, 10, 20)
	wStree2 := randChars(s.R, s.Buf, 10, 20)
	wCity := randChars(s.R, s.Buf, 10, 20)
	wState := randState(s.R, s.Buf)
	wZip := randZip(s.R, s.Buf)
	wTax := randTax(s.R)
	wYtd := 300000.00

	v := fmt.Sprintf(`(%d, '%s', '%s', '%s', '%s', '%s', '%s', %f, %f)`,
		warehouse, wName, wStree1, wStree2, wCity, wState, wZip, wTax, wYtd)
	l.InsertValue(ctx, v)

	return l.Flush(ctx)
}

func (w *Workloader) loadStock(ctx context.Context, tableID int, warehouse int) error {
	s := w.base.GetState(ctx)

	hint := fmt.Sprintf(`INSERT INTO stock%d (s_i_id, s_w_id, s_quantity, 
s_dist_01, s_dist_02, s_dist_03, s_dist_04, s_dist_05, s_dist_06, 
s_dist_07, s_dist_08, s_dist_09, s_dist_10, s_ytd, s_order_cnt, s_remote_cnt, s_data) VALUES `, tableID)

	l := load.NewBatchLoader(s.Conn, hint)

	for i := 0; i < stockPerWarehouse; i++ {
		s.Buf.Reset()
		sIID := i + 1
		sWID := warehouse
		sQuantity := randInt(s.R, 10, 100)
		sDist01 := randLetters(s.R, s.Buf, 24, 24)
		sDist02 := randLetters(s.R, s.Buf, 24, 24)
		sDist03 := randLetters(s.R, s.Buf, 24, 24)
		sDist04 := randLetters(s.R, s.Buf, 24, 24)
		sDist05 := randLetters(s.R, s.Buf, 24, 24)
		sDist06 := randLetters(s.R, s.Buf, 24, 24)
		sDist07 := randLetters(s.R, s.Buf, 24, 24)
		sDist08 := randLetters(s.R, s.Buf, 24, 24)
		sDist09 := randLetters(s.R, s.Buf, 24, 24)
		sDist10 := randLetters(s.R, s.Buf, 24, 24)
		sYtd := 0
		sOrderCnt := 0
		sRemoteCnt := 0
		sData := randOriginalString(s.R, s.Buf)

		v := fmt.Sprintf(`(%d, %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %d, %d, %d, '%s')`,
			sIID, sWID, sQuantity, sDist01, sDist02, sDist03, sDist04, sDist05, sDist06, sDist07, sDist08, sDist09, sDist10, sYtd, sOrderCnt, sRemoteCnt, sData)
		if err := l.InsertValue(ctx, v); err != nil {
			return err
		}
	}
	return l.Flush(ctx)
}

func (w *Workloader) loadDistrict(ctx context.Context, tableID int, warehouse int) error {
	s := w.base.GetState(ctx)

	hint := fmt.Sprintf(`INSERT INTO district%d (d_id, d_w_id, d_name, d_street_1, d_street_2, 
d_city, d_state, d_zip, d_tax, d_ytd, d_next_o_id) VALUES `, tableID)

	l := load.NewBatchLoader(s.Conn, hint)

	for i := 0; i < districtPerWarehouse; i++ {
		dID := i + 1
		dWID := warehouse
		dName := randChars(s.R, s.Buf, 6, 10)
		dStreet1 := randChars(s.R, s.Buf, 10, 20)
		dStreet2 := randChars(s.R, s.Buf, 10, 20)
		dCity := randChars(s.R, s.Buf, 10, 20)
		dState := randState(s.R, s.Buf)
		dZip := randZip(s.R, s.Buf)
		dTax := randTax(s.R)
		dYtd := 300000.00
		dNextOID := 3001

		v := fmt.Sprintf(`(%d, %d, '%s', '%s', '%s', '%s', '%s', '%s', %f, %f, %d)`, dID, dWID,
			dName, dStreet1, dStreet2, dCity, dState, dZip, dTax, dYtd, dNextOID)

		if err := l.InsertValue(ctx, v); err != nil {
			return err
		}
	}
	return l.Flush(ctx)
}

func (w *Workloader) loadCustomer(ctx context.Context, tableID int, warehouse int, district int) error {
	s := w.base.GetState(ctx)

	hint := fmt.Sprintf(`INSERT INTO customer%d (c_id, c_d_id, c_w_id, c_last, c_middle, c_first, 
c_street_1, c_street_2, c_city, c_state, c_zip, c_phone, c_since, c_credit, c_redit_limt,
c_discount, c_balance, c_ytd_payment, c_payment_cnt, c_delivery_cnt, c_data) VALUES `, tableID)

	l := load.NewBatchLoader(s.Conn, hint)

	for i := 0; i < customerPerDistrict; i++ {
		s.Buf.Reset()

		cID := i + 1
		cDID := district
		cWID := warehouse
		var cLast string
		if i < 1000 {
			cLast = randCLastSyllables(i, s.Buf)
		} else {
			cLast = randCLast(s.R, s.Buf)
		}
		cMiddle := "OE"
		cFirst := randChars(s.R, s.Buf, 8, 16)
		cStreet1 := randChars(s.R, s.Buf, 10, 20)
		cStreet2 := randChars(s.R, s.Buf, 10, 20)
		cCity := randChars(s.R, s.Buf, 10, 20)
		cState := randState(s.R, s.Buf)
		cZip := randZip(s.R, s.Buf)
		cPhone := randNumbers(s.R, s.Buf, 16, 16)
		cSince := time.Now().Format("2006-01-02 15:04:05")
		cCredit := "GC"
		if s.R.Intn(10) == 0 {
			cCredit = "BC"
		}
		cCreditLim := 50000.00
		cDisCount := float64(randInt(s.R, 0, 5000)) / float64(10000.0)
		cBalance := -10.00
		cYtdPayment := 10.00
		cPaymentCnt := 1
		cDeliveryCnt := 0
		cData := randChars(s.R, s.Buf, 300, 500)

		v := fmt.Sprintf(`(%d, %d, %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %f, %f, %f, %f, %d, %d, '%s')`,
			cID, cDID, cWID, cLast, cMiddle, cFirst, cStreet1, cStreet2, cCity, cState,
			cZip, cPhone, cSince, cCredit, cCreditLim, cDisCount, cBalance,
			cYtdPayment, cPaymentCnt, cDeliveryCnt, cData)
		if err := l.InsertValue(ctx, v); err != nil {
			return err
		}
	}

	return l.Flush(ctx)
}

func (w *Workloader) loadHistory(ctx context.Context, tableID int, warehouse int, district int, customer int) error {

	return nil
}

func (w *Workloader) loadOrder(ctx context.Context, tableID int, warehouse int, district int) error {
	return nil
}

func (w *Workloader) loadOrderLine(ctx context.Context, tableID int, warehouse int, district int, order int) error {
	return nil
}

func (w *Workloader) loadNewOrder(ctx context.Context, tableID int, warehouse int, district int) error {
	return nil
}