; disk2.asm
;
; This is the DISASSEMBLED source code for the Disk II controller ROM.
; It adds up to 256 bytes of program code, which is all any peripheral
; card was afforded.
;
; NOTE THAT THIS SOURCE CODE IS NOT ORIGINAL TO APPLE. I translated by
; hand from the machine code in the ROM. Any comments, etc. you see
; here, are from me--NOT APPLE.
;
; For details on the assembly instructions--what they mean and do--this
; is a good resource: 
; http://www.e-tradition.net/bytes/6502/6502_instruction_set.html
;
; Any number that has a $ in front of it means it's a hex number, vs.
; decimal.

; Our definitions for this little program. The EQU symbol is not a
; formal instruction understood by the 6502 CPU; it's a notation that
; simply means "equals"; e.g. GBASL equals $26.
                GBASL       EQU $26
                GBASH       EQU $27
                BAS2H       EQU $2B
                A1L         EQU $3C
                A1H         EQU $3D
                A3L         EQU $40
                A3H         EQU $41
                STACK       EQU $0100
                GCRALT      EQU $02D6
                XORSAV      EQU $0300
                GCRTAB      EQU $0356
                ENTRY       EQU $0800
                PHASEOFF    EQU $C080
                PHASEON     EQU $C081
                TURNON      EQU $C089
                SLCTD1      EQU $C08A
                READ        EQU $C08C
                SETRD       EQU $C08E
                WAIT        EQU $FCA8
                FINDSLOT    EQU $FF58

; This is not needed by the disk controller program itself, but is used
; by the system ROM to determine if there is a valid controller program
; here. (There is!)
00:A2 20                    LDX #$20

; Here we will write the group coded recording table for our decode
; process, which is in the GCRTAB address.
02:A0 00                    LDY #$00
04:A2 03                    LDX #$03        ; our loop begins at $03
06:86 3C        NEXTGCR     STX A1L         ; A1L tracks the loop counter
08:8A                       TXA
09:0A                       ASL A
0A:24 3C                    BIT A1L         ; if X and X<<1 have no bits in common
0C:F0 10                    BEQ CONTGCR     ; then X will not be written into GCRTAB

; This sequence of operations will prime A for the VALIDENT check. Any
; valid X register value (from which A is derived) will be one where
; these three operations results in $40; after we loop on VALIDENT and
; shift A to the right a bunch of times, we'll end up with $00 and exit
; the loop without tripping the BCS (because none of the first 6 bits
; were ever high).
0E:05 3C                    ORA A1L
10:49 FF                    EOR #$FF
12:29 7E                    AND #$7E
14:B0 08        VALIDENT    BCS CONTGCR
16:4A                       LSR A
17:D0 FB                    BNE VALIDENT

; Only certain X register values will be written via the STA
; instruction; what we do write is an iteration of Y from $00..$3F.
19:98                       TYA
1A:9D 56 03                 STA GCRTAB,X
1D:C8                       INY
1E:E8           CONTGCR     INX
1F:10 E5                    BPL NEXTGCR

; All this wrangling is here to make a record of the slot number.
; Because JSR will push the calling address into the stack, we can find
; the MSB of that address with the LDA STACK,X instruction. All of the
; ASLs will essentially push the 7 one hex digit over, so $C7 becomes
; $70. And we store that in BAS2H so we can use it to run operations on
; the peripheral.
21:20 58 FF                 JSR FINDSLOT
24:BA                       TSX
25:BD 00 01                 LDA STACK,X     ; this will load $C7 into A
28:0A                       ASL
29:0A                       ASL
2A:0A                       ASL
2B:0A                       ASL             ; and now we have $70
2C:85 2B                    STA BAS2H

; Ok, with that done, we're going to get everything set up to copy the
; zero track into RAM. NOTE: I'm not entirely sure why we're doing a
; READ from the drive before we know we have drive 1 selected and turned
; on.
2E:AA                       TAX
2F:BD 8E C0                 LDA SETRD,X
32:BD 8C C0                 LDA READ,X
35:BD 8A C0                 LDA SLCTD1,X
38:BD 89 C0                 LDA TURNON,X

; This loop is going to go through the stepper motor phases, flipping
; them off and on again. To begin with, X is $70, so we're going to work
; with phase 0 at the start.
3B:A0 50                    LDY #$50
3D:BD 80 C0     PHASELOOP   LDA PHASEOFF,X
40:98                       TYA
41:29 03                    AND #$03            ; drop all but the first 2 bits
43:0A                       ASL                 ; and shift over
44:05 2B                    ORA BAS2H           ; and add that to $70
46:AA                       TAX
47:BD 81 C0                 LDA PHASEON,X
4A:A9 56                    LDA #$56

; In at least one implementation (notably WinApple), the opcode below is
; rewritten as `A9 00 EA`, which is equivalent to:
;   LDA #$00
;   NOP
; This would essentially remove the WAIT call. The WAIT subroutine will,
; in the course of its operation, leave $00 in A, which explains the LDA
; #$00 opcode sequence. The NOP is there to replace the third byte
; (which was part of the JSR address in its original form).
4C:20 A8 FC                 JSR WAIT            ; wait for the motor
4F:88                       DEY
50:10 EB                    BPL PHASELOOP

; We're setting things up so we can start writing our decoded data into
; the $08 page in memory, which is where we will ultimately jump to once
; we finish going through track zero.
52:85 26                    STA GBASL           ; A is $00 by this point
54:85 3D                    STA A1H
56:85 41                    STA A3H
58:A9 08                    LDA #$08
5A:85 27                    STA GBASH           ; so GBASH/L will hold $0800

; We're going to check to see if we are at a header marker.
5C:18           CHKHD       CLC
5D:08           CHKHDC      PHP                 ; hang onto the status for later

; Read byte from the disk (BPL is used here because anything that
; doesn't have bit 7 high is bad data in 6-and-2 encoding).
5E:BD 8C C0     READHD1     LDA READ,X
61:10 FB                    BPL READHD1
63:49 D5        CHKHD1      EOR #$D5
65:D0 F7                    BNE READHD1         ; try again

; Look for the second header byte
67:BD 8C C0     READHD2     LDA READ,X
6A:10 FB                    BPL READHD2
6C:C9 AA        CHKHD2      CMP #$AA
6E:D0 F3                    BNE CHKHD1
70:EA                       NOP                 ; I don't know why we NOP here

; Third header byte
71:BD 8C C0     READHD3     LDA READ,X
74:10 FB                    BPL READHD3
76:C9 96                    CMP #$96            ; is this the end of a track marker?
78:F0 09                    BEQ METADATA        ; seems to be!
7A:28                       PLP
7B:90 DF                    BCC CHKHD           ; if A < $96, keep seeking for a header byte 
7D:49 AD                    EOR #$AD            ; if NOT, then this might be the end of a sector header
7F:F0 25                    BEQ DECODE          ; so let's get decoding!
81:D0 D9                    BNE CHKHD           ; Some other byte we didn't expect...

; The metadata is 4-and-4 encoded, which are two bytes that are read in
; sequence and then AND'd together. The second in the sequence is what
; will stay behind in A3L; we'll read 3 sequences in all.
83:A0 03        METADATA    LDY #$03
85:85 40        AGAIN44     STA A3L
87:BD 8C C0     FIRST44     LDA READ,X          ; read a byte
8A:10 FB                    BPL FIRST44
8C:2A                       ROL A
8D:85 3C                    STA A1L
8F:BD 8C C0     SECOND44    LDA READ,X          ; read another byte
92:10 FB                    BPL SECOND44
94:25 3C                    AND A1L             ; intersect with the shifted FIRST44
96:88                       DEY
97:D0 EC                    BNE AGAIN44

; This is going to pull from before we began checking for a header
99:28                       PLP
9A:C5 3D                    CMP A1H
9C:D0 BE                    BNE CHKHD

; A3H can only be $00, and A3L will have been $96 from the last header
; byte we read; since $96 - $00 will of course not be zero, this will
; force a branch back to CHKHD. Why we have this code here is unclear to
; me.
9E:A5 40                    LDA A3L
A0:C5 41                    CMP A3H
A2:D0 B8                    BNE CHKHD

; If C is set, we will jump back to read the next header, _but_ we will
; not execute the CLC instruction.
A4:B0 B7                    BCS CHKHDC

; As we decode bytes, we're referencing the GCRTAB entries we built
; earlier but from a slightly different address point (hence GCRALT).
; But make no mistake--we're EORing with GCRTAB data. Note that XORSAV
; is an entry point ($0300) which is conveniently(!) $56 less than
; GCRTAB ($0356). It is, though, really just a place to stash the
; intermediate data.
A6:A0 56        DECODE      LDY #$56        ; loop this many times...
A8:84 3C        SAV2BITS    STY A1L         ; save in A1L, because we use Y to read
AA:BC 8C C0     DECBYTE2    LDY READ,X
AD:10 FB                    BPL DECBYTE
AF:59 D6 02                 EOR GCRALT,Y
B2:A4 3C                    LDY A1L
B4:88                       DEY             ; decrement the loop counter
B5:99 00 03                 STA XORSAV,Y    ; hang onto the EOR data
B8:D0 EE                    BNE SAV2BITS

; Looping from zero, now, we're going to write all that intermediate
; data into the $0800 page (which is what (GBASL),Y resolves to),
; counting up from $0800..$08FF.
BA:84 3C        SAV6BITS    STY A1L
BC:BC 8C C0     DECBYTE6    LDY READ,X
BF:10 FB                    BPL DECBYTE6
C1:59 D6 02                 EOR GCRALT,Y
C4:A4 3C                    LDY A1L
C6:91 26                    STA (GBASL),Y
C8:C8                       INY
C9:D0 EF                    BNE SAV6BITS

; We read ONE more byte, then determine if we need to check for another
; header again.
CB:BC 8C C0     FINBYTE     LDY READ,X
CE:10 FB                    BPL FINBYTE
D0:59 D6 02                 EOR XORTMP,Y

; We may also get to here because it looks like we didn't write data
; properly into the ENTRY page.
D3:D0 87        BADDATA     BNE CHKHD       ; another sector?

; We're using the first 89 ($56) bytes in XORSAV (which, remember, is
; $56 less than the GCRTAB address point); if we go below $00 (rolling
; over to $FF), start over.
;
; The bytes we've already written into the ENTRY page need to have those
; 2 bits we compiled into those 89 bytes pushed back into the data.
D5:A0 00                    LDY #$00
D7:A2 56        BITLOOP     LDX #$56
D9:CA           WRITELOOP   DEX
DA:30 FB                    BMI BITLOOP     ; start over if we went past $00
DC:B1 26                    LDA (GBASL),Y   ; load $0800 + Y
DE:5E 00 03                 LSR XORSAV,X    ; move bit 0 into carry
E1:2A                       ROL A           ; now load carry into A, plus the orig contents
E2:5E 00 03                 LSR XORSAV,X    ; shift the former bit 1 (now bit 0) into carry again
E5:2A                       ROL A           ; and again load into A; now we have all 8 bits
E6:91 26                    STA (GBASL),Y   ; and save it back to $0800 + Y
E8:C8                       INY
E9:D0 EE                    BNE WRITELOOP   ; we'll loop here 256 times

; We're in the home stretch... we're just double-checking if we copied
; things into the ENTRY page properly.
EB:E6 27                    INC GBASH       ; so now GBASL/H is $0900
ED:E6 3D                    INC A1H
EF:A5 3D                    LDA A1H
F1:CD 00 08                 CMP ENTRY       ; if A < ENTRY
F4:A6 2B                    LDX BAS2H
F6:90 DB                    BCC BADDATA     ; then go back and try again
F8:4C 01 08                 JMP ENTRY+1     ; otherwise, let's boot the software!
