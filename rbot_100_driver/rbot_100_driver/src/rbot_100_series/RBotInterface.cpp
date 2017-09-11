/*********************************************************************
*
* Software License Agreement (BSD License)
*
*  Copyright (c) 2010, ISR University of Coimbra.
*  All rights reserved.
*
*  Redistribution and use in source and binary forms, with or without
*  modification, are permitted provided that the following conditions
*  are met:
*
*   * Redistributions of source code must retain the above copyright
*     notice, this list of conditions and the following disclaimer.
*   * Redistributions in binary form must reproduce the above
*     copyright notice, this list of conditions and the following
*     disclaimer in the documentation and/or other materials provided
*     with the distribution.
*   * Neither the name of the ISR University of Coimbra nor the names of its
*     contributors may be used to endorse or promote products derived
*     from this software without specific prior written permission.
*
*  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
*  "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
*  LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
*  FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
*  COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
*  INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
*  BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
*  LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
*  CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
*  LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
*  ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
*  POSSIBILITY OF SUCH DAMAGE.
*
* Author: Gonçalo Cabrita on 19/05/2010
*********************************************************************/
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <stdio.h>
#include <string>
#include <netinet/in.h>
#include <sys/types.h>

#include "rbot_100_series/RBotInterface.h"

// *****************************************************************************
// Constructor
rbot::RBotInterface::RBotInterface(const char * new_serial_port)
{	
	port_name_ = new_serial_port;

	OImode_ = OI_MODE_OFF;
	
	this->resetOdometry();
	
	encoder_counts_[LEFT] = -1;
	encoder_counts_[RIGHT] = -1;
	
	last_encoder_counts_[LEFT] = 0;
	last_encoder_counts_[RIGHT] = 0;
	
	num_of_packets_ = 0;
	sensor_packets_ = NULL;
	packets_size_ = 0;
	
	// Default packets
	OI_Packet_ID default_packets[2] = {OI_PACKET_RIGHT_ENCODER, OI_PACKET_LEFT_ENCODER};
	this->setSensorPackets(default_packets, 2, OI_PACKET_RIGHT_ENCODER_SIZE + OI_PACKET_LEFT_ENCODER_SIZE);

	serial_port_ = new cereal::CerealPort();
}


// *****************************************************************************
// Destructor
rbot::RBotInterface::~RBotInterface()
{
	// Clean up!
	delete serial_port_;
}


// *****************************************************************************
// Open the serial port
int rbot::RBotInterface::openSerialPort(bool full_init)
{
	try{ serial_port_->open(port_name_.c_str(), SERIAL_SPEED); }
	catch(cereal::Exception& e){ return(-1); }

	this->startOI(full_control);

	return(0);
}


// *****************************************************************************
// Set the mode
int rbot::RBotInterface::startOI(bool full_init)
{	
	char buffer[1];

	usleep(OI_DELAY_MODECHANGE_MS * 1e3);
	buffer[0] = (char)OI_OPCODE_START;
	try{ serial_port_->write(buffer, 1); }
	catch(cereal::Exception& e){ return(-1); }
	OImode_ = OI_MODE_PASSIVE;

	usleep(OI_DELAY_MODECHANGE_MS * 1e3);
	buffer[0] = (char)OI_OPCODE_CONTROL;
	try{ serial_port_->write(buffer, 1); }
	catch(cereal::Exception& e){ return(-1); }
	OImode_ = OI_MODE_SAFE;
	
	if(full_control)
	{
		usleep(OI_DELAY_MODECHANGE_MS * 1e3);
		buffer[0] = (char)OI_OPCODE_FULL;
		try{ serial_port_->write(buffer, 1); }
		catch(cereal::Exception& e){ return(-1); }
		OImode_ = OI_MODE_FULL;
	}
	return(0);
}


bool sendCommand(char* command)
{
	unsigned char cmd_buffer[4];
	cmd_buffer[0] = (unsigned char)OI_OPCODE_LEDS;
	cmd_buffer[1] = debris | spot<<1 | dock<<2 | check_robot<<3;
	cmd_buffer[2] = power_color;
	cmd_buffer[3] = power_intensity;
	
	try{ serial_port_->write((char*)cmd_buffer, 4); }
	catch(cereal::Exception& e){ return(-1); }
	return(0);

}


// *****************************************************************************
// Close the serial port
int rbot::RBotInterface::closeSerialPort()
{
	this->drive(0.0, 0.0);
	usleep(OI_DELAY_MODECHANGE_MS * 1e3);

	try{ serial_port_->close(); }
	catch(cereal::Exception& e){ return(-1); }

	return(0);
}


// *****************************************************************************
// Send an OP code to the rbot
int rbot::RBotInterface::sendOpcode(OI_Opcode code)
{
	char to_send = code;
	try{ serial_port_->write(&to_send, 1); }
	catch(cereal::Exception& e){ return(-1); }
	return(0);
}


// *****************************************************************************
// Power down the rbot
int rbot::RBotInterface::powerDown()
{
	return sendOpcode(OI_OPCODE_POWER);
}


// *****************************************************************************
// Set the speeds
int rbot::RBotInterface::drive(double linear_speed, double angular_speed)
{
	int left_speed_mm_s = (int)((linear_speed - ROOMBA_AXLE_LENGTH * angular_speed / 2) * 1e3);		// Left wheel velocity in mm/s
	int right_speed_mm_s = (int)((linear_speed + ROOMBA_AXLE_LENGTH * angular_speed / 2) * 1e3);	// Right wheel velocity in mm/s
	
	return this->driveDirect(left_speed_mm_s, right_speed_mm_s);
}


// *****************************************************************************
// Set the motor speeds
/*int rbot::RBotInterface::driveDirect(int left_speed, int right_speed)
{
	// Limit velocity
	int16_t left_speed_mm_s = MAX(left_speed, -ROOMBA_MAX_LIN_VEL_MM_S);
	left_speed_mm_s = MIN(left_speed, ROOMBA_MAX_LIN_VEL_MM_S);
	int16_t right_speed_mm_s = MAX(right_speed, -ROOMBA_MAX_LIN_VEL_MM_S);
	right_speed_mm_s = MIN(right_speed, ROOMBA_MAX_LIN_VEL_MM_S);
	
	// Compose comand
	char cmd_buffer[5];
	cmd_buffer[0] = (char)OI_OPCODE_DRIVE_DIRECT;
	cmd_buffer[1] = (char)(right_speed_mm_s >> 8);
	cmd_buffer[2] = (char)(right_speed_mm_s & 0xFF);
	cmd_buffer[3] = (char)(left_speed_mm_s >> 8);
	cmd_buffer[4] = (char)(left_speed_mm_s & 0xFF);

	try{ serial_port_->write(cmd_buffer, 5); }
	catch(cereal::Exception& e){ return(-1); }

	return(0);
}*/
int rbot::RBotInterface::driveDirect(int left_speed, int right_speed)
{
	// Limit velocity
	int16_t left_speed_mm_s = MAX(left_speed, -ROOMBA_MAX_LIN_VEL_MM_S);
	left_speed_mm_s = MIN(left_speed, ROOMBA_MAX_LIN_VEL_MM_S);
	int16_t right_speed_mm_s = MAX(right_speed, -ROOMBA_MAX_LIN_VEL_MM_S);
	right_speed_mm_s = MIN(right_speed, ROOMBA_MAX_LIN_VEL_MM_S);
	
	// Compose comand
	char cmd_buffer[5];
	cmd_buffer[0] = (char)OI_OPCODE_DRIVE_DIRECT;
	cmd_buffer[1] = (char)(right_speed_mm_s >> 8);
	cmd_buffer[2] = (char)(right_speed_mm_s & 0xFF);
	cmd_buffer[3] = (char)(left_speed_mm_s >> 8);
	cmd_buffer[4] = (char)(left_speed_mm_s & 0xFF);

	try{ serial_port_->write(cmd_buffer, 5); }
	catch(cereal::Exception& e){ return(-1); }

	return(0);
}


// *****************************************************************************
// Set the sensors to read
int rbot::RBotInterface::setSensorPackets(OI_Packet_ID * new_sensor_packets, int new_num_of_packets, size_t new_buffer_size)
{
	if(sensor_packets_ == NULL)
	{
		delete [] sensor_packets_;
	}
	
	num_of_packets_ = new_num_of_packets;
	sensor_packets_ = new OI_Packet_ID[num_of_packets_];
	
	for(int i=0 ; i<num_of_packets_ ; i++)
	{
		sensor_packets_[i] = new_sensor_packets[i];
	}

	stream_defined_ = false;
	packets_size_ = new_buffer_size;
	return(0);
}


// *****************************************************************************
// Read the sensors
int rbot::RBotInterface::getSensorPackets(int timeout)
{
	char cmd_buffer[num_of_packets_+2];
	char data_buffer[packets_size_];

	// Fill in the command buffer to send to the robot
	cmd_buffer[0] = (char)OI_OPCODE_QUERY;			// Query
	cmd_buffer[1] = num_of_packets_;				// Number of packets
	for(int i=0 ; i<num_of_packets_ ; i++)
	{
		cmd_buffer[i+2] = sensor_packets_[i];		// The packet IDs
	}

	try{ serial_port_->write(cmd_buffer, num_of_packets_+2); }
	catch(cereal::Exception& e){ return(-1); }
	
	try{ serial_port_->readBytes(data_buffer, packets_size_, timeout); }
	catch(cereal::Exception& e){ return(-1); }
	
	return this->parseSensorPackets((unsigned char*)data_buffer, packets_size_);
}


// *****************************************************************************
// Parse sensor data
int rbot::RBotInterface::parseSensorPackets(unsigned char * buffer , size_t buffer_lenght)
{	
	if(buffer_lenght != packets_size_)
	{
		// Error wrong packet size
		return(-1);
	}

	int i = 0;
	unsigned int index = 0;
	while(index < packets_size_)
	{
		if(sensor_packets_[i]==31)		// OI_PACKET_BUMPS_DROPS
		{
			index += parseBumpersAndWheeldrops(buffer, index);
			i++;
		}
		if(sensor_packets_[i]==OI_PACKET_WALL)
		{
			index += parseWall(buffer, index);
			i++;
		}

	}
	return(0);
}

int rbot::RBotInterface::parseBumpersAndWheeldrops(unsigned char * buffer, int index)
{
	// Bumps, wheeldrops	
	this->bumper_[RIGHT] = (buffer[index]) & 0x01;
	this->bumper_[LEFT] = (buffer[index] >> 1) & 0x01;
	this->wheel_drop_[RIGHT] = (buffer[index] >> 2) & 0x01;
	this->wheel_drop_[LEFT] = (buffer[index] >> 3) & 0x01;
	
	return 2;	// OI_PACKET_BUMPS_DROPS_SIZE
}

int rbot::RBotInterface::parseWall(unsigned char * buffer, int index)
{
	// Wall
	this->wall_ = buffer[index] & 0x01;
	
	return 2;	// OI_PACKET_WALL_SIZE
}
	
int rbot::RBotInterface::parseDistance(unsigned char * buffer, int index)
{
	// Distance
	this->distance_ = buffer2signed_int(buffer, index);
	
	return 2;	//OI_PACKET_DISTANCE_SIZE
}

int rbot::RBotInterface::parseAngle(unsigned char * buffer, int index)
{
	// Angle
	this->angle_ = buffer2signed_int(buffer, index);

	return 2;	//OI_PACKET_ANGLE_SIZE
}
	
int rbot::RBotInterface::parseChargingState(unsigned char * buffer, int index)
{
	// Charging State
	unsigned char byte = buffer[index];
	
	this->power_cord_ = (byte >> 0) & 0x01;
	this->dock_ = (byte >> 1) & 0x01;

	return 2;	// OI_PACKET_CHARGING_STATE_SIZE
}

int rbot::RBotInterface::parseVoltage(unsigned char * buffer, int index)
{
	// Voltage
	this->voltage_ = (float)(buffer2unsigned_int(buffer, index) / 1000.0);

	return 1;	// OI_PACKET_VOLTAGE_SIZE
}

int rbot::RBotInterface::parseCurrent(unsigned char * buffer, int index)
{
	// Current
	this->current_ = (float)(buffer2signed_int(buffer, index) / 1000.0);

	return 1;	// OI_PACKET_CURRENT_SIZE
}

int rbot::RBotInterface::parseTemperature(unsigned char * buffer, int index)
{
	// Temperature
	this->temperature_ = (char)(buffer[index]);

	return 1;	// OI_PACKET_TEMPERATURE_SIZE
}

int rbot::RBotInterface::parseBatteryCharge(unsigned char * buffer, int index)
{
	// Charge
	this->charge_ = (float)(buffer2unsigned_int(buffer, index) / 1000.0);

	return 1; 	//OI_PACKET_BATTERY_CHARGE_SIZE;
}

int rbot::RBotInterface::parseBatteryCapacity(unsigned char * buffer, int index)
{
	// Capacity
	this->capacity_ = (float)(buffer2unsigned_int(buffer, index) / 1000.0);

	return 1;	//OI_PACKET_BATTERY_CAPACITY_SIZE;
}
	
int rbot::RBotInterface::parseChargingSource(unsigned char * buffer, int index)
{
	// Charging soruces available
	this->power_cord_ = (buffer[index] >> 0) & 0x01;
	this->dock_ = (buffer[index] >> 1) & 0x01;

	return 1;	//OI_PACKET_CHARGE_SOURCES_SIZE;
}

int rbot::RBotInterface::parseRightEncoderCounts(unsigned char * buffer, int index)
{
	// Right encoder counts
	uint16_t right_encoder_counts = buffer2unsigned_int(buffer, index);

	//printf("Right Encoder: %d\n", rightEncoderCounts);

	if(encoder_counts_[RIGHT] == -1 || right_encoder_counts == last_encoder_counts_[RIGHT])	// First time, we need 2 to make it work!
	{
		encoder_counts_[RIGHT] = 0;
	}
	else
	{
		encoder_counts_[RIGHT] = (int)(right_encoder_counts - last_encoder_counts_[RIGHT]);
		
		if(encoder_counts_[RIGHT] > ROOMBA_MAX_ENCODER_COUNTS/10) encoder_counts_[RIGHT] = encoder_counts_[RIGHT] - ROOMBA_MAX_ENCODER_COUNTS;
		if(encoder_counts_[RIGHT] < -ROOMBA_MAX_ENCODER_COUNTS/10) encoder_counts_[RIGHT] = ROOMBA_MAX_ENCODER_COUNTS + encoder_counts_[RIGHT];
	}
	last_encoder_counts_[RIGHT] = right_encoder_counts;
	
	return 2;	//OI_PACKET_RIGHT_ENCODER_SIZE;
}

int rbot::RBotInterface::parseLeftEncoderCounts(unsigned char * buffer, int index)
{
	// Left encoder counts
	uint16_t left_encoder_counts = buffer2unsigned_int(buffer, index);

	//printf("Left Encoder: %d\n", leftEncoderCounts);

	if(encoder_counts_[LEFT] == -1 || left_encoder_counts == last_encoder_counts_[LEFT])	// First time, we need 2 to make it work!
	{
		encoder_counts_[LEFT] = 0;
	}
	else
	{
		encoder_counts_[LEFT] = (int)(left_encoder_counts - last_encoder_counts_[LEFT]);
		
		if(encoder_counts_[LEFT] > ROOMBA_MAX_ENCODER_COUNTS/10) encoder_counts_[LEFT] = encoder_counts_[LEFT] - ROOMBA_MAX_ENCODER_COUNTS;
		if(encoder_counts_[LEFT] < -ROOMBA_MAX_ENCODER_COUNTS/10) encoder_counts_[LEFT] = ROOMBA_MAX_ENCODER_COUNTS + encoder_counts_[LEFT];
	}
	last_encoder_counts_[LEFT] = left_encoder_counts;
	
	return 2;	//OI_PACKET_LEFT_ENCODER_SIZE;
}
	
int rbot::RBotInterface::parseLeftMotorCurrent(unsigned char * buffer, int index)
{
	// Left motor current
	this->motor_current_[LEFT] = buffer2signed_int(buffer, index);
	
	return 2;	//OI_PACKET_LEFT_MOTOR_CURRENT_SIZE;
}

int rbot::RBotInterface::parseRightMotorCurrent(unsigned char * buffer, int index)
{
	// Left motor current
	this->motor_current_[RIGHT] = buffer2signed_int(buffer, index);
	
	return 2;	//OI_PACKET_RIGHT_MOTOR_CURRENT_SIZE;
}

int rbot::RBotInterface::buffer2signed_int(unsigned char * buffer, int index)
{
	int16_t signed_int;
	
	memcpy(&signed_int, buffer+index, 2);
	signed_int = ntohs(signed_int);
	
	return (int)signed_int;
}

int rbot::RBotInterface::buffer2unsigned_int(unsigned char * buffer, int index)
{
	uint16_t unsigned_int;

	memcpy(&unsigned_int, buffer+index, 2);
	unsigned_int = ntohs(unsigned_int);
	
	return (int)unsigned_int;
}


// *****************************************************************************
// Calculate Roomba odometry
void rbot::RBotInterface::calculateOdometry()
{	
	double dist = (encoder_counts_[RIGHT]*ROOMBA_PULSES_TO_M + encoder_counts_[LEFT]*ROOMBA_PULSES_TO_M) / 2.0; 
	double ang = (encoder_counts_[RIGHT]*ROOMBA_PULSES_TO_M - encoder_counts_[LEFT]*ROOMBA_PULSES_TO_M) / -ROOMBA_AXLE_LENGTH;

	// Update odometry
	this->odometry_yaw_ = NORMALIZE(this->odometry_yaw_ + ang);			// rad
	this->odometry_x_ = this->odometry_x_ + dist*cos(odometry_yaw_);		// m
	this->odometry_y_ = this->odometry_y_ + dist*sin(odometry_yaw_);		// m
}


// *****************************************************************************
// Reset Roomba odometry
void rbot::RBotInterface::resetOdometry()
{
	this->setOdometry(0.0, 0.0, 0.0);
}


// *****************************************************************************
// Set Roomba odometry
void rbot::RBotInterface::setOdometry(double new_x, double new_y, double new_yaw)
{
	this->odometry_x_ = new_x;
	this->odometry_y_ = new_y;
	this->odometry_yaw_ = new_yaw;
}


// *****************************************************************************
// Go to the dock
int rbot::RBotInterface::goDock()
{
	return sendOpcode(OI_OPCODE_FORCE_DOCK);
}


// EOF
