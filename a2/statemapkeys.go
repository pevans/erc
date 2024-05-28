package a2

import "fmt"

const (
	kbLastKey           = 100
	kbStrobe            = 101
	kbKeyDown           = 102
	memRead             = 200
	memWrite            = 201
	memReadSegment      = 202
	memWriteSegment     = 203
	memAuxSegment       = 204
	memMainSegment      = 205
	pcExpansion         = 300
	pcSlotC3            = 301
	pcSlotCX            = 302
	pcExpSlot           = 303
	pcROMSegment        = 304
	bankRead            = 401
	bankWrite           = 402
	bankDFBlock         = 403
	bankSysBlock        = 404
	bankReadAttempts    = 406
	bankSysBlockSegment = 407
	bankROMSegment      = 408
	displayAltChar      = 500
	displayCol80        = 501
	displayStore80      = 502
	displayPage2        = 503
	displayText         = 504
	displayMixed        = 505
	displayHires        = 506
	displayIou          = 507
	displayDoubleHigh   = 508
	displayRedraw       = 509
	displayAuxSegment   = 510
	diskComputer        = 600
)

func stateMapKeyToString(key any) string {
	intKey, ok := key.(int)
	if !ok {
		return fmt.Sprintf("unknown non-int key (%v)", key)
	}

	switch intKey {
	case bankRead:
		return "bankRead"
	case bankWrite:
		return "bankWrite"
	case bankDFBlock:
		return "bankDFBlock"
	case bankSysBlock:
		return "bankSysBlock"
	case bankReadAttempts:
		return "bankReadAttempts"
	case bankSysBlockSegment:
		return "bankSysBlockSegment"
	case bankROMSegment:
		return "bankROMSegment"
	case diskComputer:
		return "diskComputer"
	case displayAltChar:
		return "displayAltChar"
	case displayCol80:
		return "displayCol80"
	case displayStore80:
		return "displayStore80"
	case displayPage2:
		return "displayPage2"
	case displayText:
		return "displayText"
	case displayMixed:
		return "displayMixed"
	case displayHires:
		return "displayHires"
	case displayIou:
		return "displayIou"
	case displayDoubleHigh:
		return "displayDoubleHigh"
	case displayRedraw:
		return "displayRedraw"
	case displayAuxSegment:
		return "displayAuxSegment"
	case kbLastKey:
		return "kbLastkey"
	case kbStrobe:
		return "kbStrobe"
	case kbKeyDown:
		return "kbKeyDown"
	case memRead:
		return "memRead"
	case memWrite:
		return "memWrite"
	case memReadSegment:
		return "memReadSegment"
	case memWriteSegment:
		return "memWriteSegment"
	case memAuxSegment:
		return "memAuxSegment"
	case memMainSegment:
		return "memMainSegment"
	case pcExpansion:
		return "pcExpansion"
	case pcSlotC3:
		return "pcSlotC3"
	case pcSlotCX:
		return "pcSlotCX"
	case pcExpSlot:
		return "pcExpSlot"
	case pcROMSegment:
		return "pcROMSegment"
	}

	return fmt.Sprintf("unknown (%v)", key)
}
