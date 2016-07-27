#ifndef __ISRS_H__
#define __ISRS_H__
/* This defines what the stack looks like after an ISR was running */
struct regs
{
    unsigned int gs, fs, es, ds;      /* pushed the segs last */
    unsigned int edi, esi, ebp, esp, ebx, edx, ecx, eax;  /* pushed by 'pusha' */
    unsigned int int_no, err_code;    /* our 'push byte #' and ecodes do this */
    unsigned int eip, cs, eflags, useresp, ss;   /* pushed by the processor automatically */ 
};
extern void idt_set_gate(unsigned char num, unsigned long base, unsigned short sel, unsigned char flags);
extern void terminal_writestring(const char*);
#endif