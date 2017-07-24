#!/usr/bin/env python
#import re
import serial
import time
import pygame, sys
from pygame.locals import *

pygame.init()
pygame.display.set_mode((100,100))

ser = None
initstr = [0xFE, 0x20, 0x01, 0x3F, 0x5E]
init = [[0xFE, 0x08, 0x00, 0x08], 
       [0x00, 0x08, 0x00, 0x08], 
     #  [0x00, 0x16, 0x00, 0x16], 
       [0x01, 0x08, 0x00, 0x09],
       [0x02, 0x08, 0x00, 0x0A], 
       [0x05, 0x27, 0x01, 0xF4, 0x21], 
       [0x05, 0x72, 0x04, 0x7F,  0x7F,  0x7F,  0x7F, 0x77]]
c_back = "00 D2 04 05 00 FB FF D5" # looks like "go back"
c_stop = "00 D2 04 00 00 00 00 D6" # looks like "stop"
c_fwd = "00 D2 04 FB FF 05 00 D5" # looks like "go fwd"
c_left = "00 D2 04 FB FF FB FF CA" # looks like "left/counterclockwise"
c_right = "00 D2 04 05 00 05 00 E0" # looks like "right/clockwise"
c_hr = "05 7C 03 02 00 00 86"    # looks like "head right" movement
c_hl = "05 7C 03 FE 00 00 82"    # looks like "head left" movement
c_hd = "05 7C 03 00 02 00 86"    # "head down"
c_hu = "05 7C 03 00 FE 00 82"    # "head up"
c_unlock = "FE 20 01 FF 1E"      # Unlock Live mode

def send_str(s):
  return [chr(int(i, 16)) for i in s.split(" ")]

def send_arr(ar):
  if len(ar) < 1:
    return
  elif type(ar[0]) == int:
    ar = [chr(i) for i in ar]
  return [chr(0x02)] + ar + [chr(0x03)]

def send(s, comment=None):
  ''' Sends a string to port. Takes string of codes, array of codes and exact string '''
  if comment == None:
    print comment 
  if type(s) == str:
    if ' ' in s:
      s = send_str(s)
  elif type(s) == list:
    pass
  else:
    raise("Wrong format for command!")
  if not s[0] == chr(0x02):
    s = send_arr(s)
  #try:
  ser.write("".join(s))
  #except:
  #  print ("Serial port issue!")
  #  pass

if __name__ == '__main__':
  try:
    ser = serial.Serial('/dev/tty.SLAB_USBtoUART', 57600, timeout=0.5)
  except:
    print ("Serial init failed!")
    ser = None
  # Initialization
  for i in range(1, 10):
    send(initstr)
  print("Init done")
  for s in init:
    send(s)
    time.sleep(0.3)
  #msg = ser.read(ser.inWaiting())
  #if len(msg) > 0:
  #    print(("Received after query: ") + " ".join([hex(ord(x)) for x in msg]) + ":    " + msg)
  # Control loop
  while True:
    for event in pygame.event.get():
      if event.type == QUIT: sys.exit()
      #if event.type == KEYDOWN:
      #  print("Pressed %d" % event.dict['key'])
      if event.type == KEYDOWN and event.dict['key'] == ord('x'):
         send(c_fwd, 'Forward')
      if event.type == KEYDOWN and event.dict['key'] == ord('w'):
         send(c_back, 'Backward')
      if event.type == KEYDOWN and event.dict['key'] == ord('a'):
         send(c_left, 'Left')
      if event.type == KEYDOWN and event.dict['key'] == ord('d'):
         send(c_right, 'Right')
      if event.type == KEYDOWN and event.dict['key'] == ord('s'):
         send(c_stop, 'Stop')
      if event.type == KEYDOWN and event.dict['key'] == ord('u'):
         send(c_unlock, 'Unlock')
      if event.type == KEYDOWN and event.dict['key'] == 273:      # up
         send(c_hu, 'Head up')
      if event.type == KEYDOWN and event.dict['key'] == 274:      # down
         send(c_hd, 'Head down')
      if event.type == KEYDOWN and event.dict['key'] == 275:      # right
         send(c_hr, 'Head right')
      if event.type == KEYDOWN and event.dict['key'] == 276:      # left
         send(c_hl, 'Head left')
    pygame.event.pump()

  ser.close()
