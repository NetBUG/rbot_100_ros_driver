#include "cereal_port/CerealPort.h"

// OI Modes
#define OI_MODE_OFF				0
#define OI_MODE_PASSIVE			1
#define OI_MODE_SAFE			2
#define OI_MODE_FULL			3

// Delay after mode change in ms
#define OI_DELAY_MODECHANGE_MS	20

// Charging states
#define OI_CHARGING_NO			0
#define OI_CHARGING_RECOVERY	1
#define OI_CHARGING_CHARGING	2
#define OI_CHARGING_TRICKLE		3
#define OI_CHARGING_WAITING		4
#define OI_CHARGING_ERROR		5

// IR Characters
#define FORCE_FIELD						161
#define GREEN_BUOY						164
#define GREEN_BUOY_FORCE_FIELD			165
#define RED_BUOY						168
#define RED_BUOY_FORCE_FIELD			169
#define RED_BUOY_GREEN_BUOY				172
#define RED_BUOY_GREEN_BUOY_FORCE_FIELD	173
#define VIRTUAL_WALL					162

// Positions
#define LEFT				0
#define RIGHT				1
#define FRONT_LEFT			2
#define FRONT_RIGHT			3
#define CENTER_LEFT			4
#define CENTER_RIGHT		5
#define OMNI				2
#define MAIN_BRUSH			2
#define SIDE_BRUSH			3

// Buttons
#define BUTTON_CLOCK		7
#define BUTTON_SCHEDULE		6
#define BUTTON_DAY			5
#define BUTTON_HOUR			4
#define BUTTON_MINUTE		3
#define BUTTON_DOCK			2
#define BUTTON_SPOT			1
#define BUTTON_CLEAN		0

// Roomba Dimensions
#define ROOMBA_BUMPER_X_OFFSET		0.050
#define ROOMBA_DIAMETER				0.330
#define ROOMBA_AXLE_LENGTH			0.235

#define ROOMBA_MAX_LIN_VEL_MM_S		500
#define ROOMBA_MAX_ANG_VEL_RAD_S	2  
#define ROOMBA_MAX_RADIUS_MM		2000

//! Roomba max encoder counts
#define ROOMBA_MAX_ENCODER_COUNTS	65535
//! Roomba encoder pulses to meter constant
#define ROOMBA_PULSES_TO_M			0.000445558279992234

#define MAX_PATH 32


unsigned char RBOT_INITSTR[5] =  		{0xFE, 0x20, 0x01, 0x3F, 0x5E};
unsigned char RBOT_INIT[4][4] =  	 	{{0xFE, 0x08, 0x00, 0x08}, 
       {0x00, 0x08, 0x00, 0x08}, 
//     #  [0x00, 0x16, 0x00, 0x16], 
       {0x01, 0x08, 0x00, 0x09},
       {0x02, 0x08, 0x00, 0x0A}}; 
unsigned char RBOT_INIT2[5] =           {0x05, 0x27, 0x01, 0xF4, 0x21};
unsigned char RBOT_INIT3[8] = 		{0x05, 0x72, 0x04, 0x7F,  0x7F,  0x7F,  0x7F, 0x77};
unsigned char RBOT_BACK[8] = {0x00, 0xD2, 0x04, 0x05, 0x00, 0xFB, 0xFF, 0xD5}; // looks like {go back"
unsigned char RBOT_STOP[8] = {0x00, 0xD2, 0x04, 0x00, 0x00, 0x00, 0x00, 0xD6}; // looks like {stop};
unsigned char RBOT_FWD[8] = {0x00, 0xD2, 0x04, 0xFB, 0xFF, 0x05, 0x00, 0xD5}; // looks like {go fwd};
unsigned char RBOT_CCW[8] = {0x00, 0xD2, 0x04, 0xFB, 0xFF, 0xFB, 0xFF, 0xCA}; // looks like {left/counterclockwise};
unsigned char RBOT_CW[8] = {0x00, 0xD2, 0x04, 0x05, 0x00, 0x05, 0x00, 0xE0}; // looks like {right/clockwise};
unsigned char RBOT_HEAD_RIGHT[8] = {0x05, 0x7C, 0x03, 0x02, 0x00, 0x00, 0x86};    // looks like {head right}; movement
unsigned char RBOT_HEAD_LEFT[8] = {0x05, 0x7C, 0x03, 0xFE, 0x00, 0x00, 0x82};    // looks like {head left}; movement
unsigned char RBOT_HEAD_DOWN[8] = {0x05, 0x7C, 0x03, 0x00, 0x02, 0x00, 0x86};    // {head down};
unsigned char RBOT_HEAD_UP[8] = {0x05, 0x7C, 0x03, 0x00, 0xFE, 0x00, 0x82};    // {head up};
unsigned char RBOT_UNLOCK[8] = {0xFE, 0x20, 0x01, 0xFF, 0x1E};      // Unlock Live mode

#define SERIAL_SPEED		57600		// RBot uses 57600, 8 bit, no parity check

#ifndef MIN
#define MIN(a,b) ((a < b) ? (a) : (b))
#endif
#ifndef MAX
#define MAX(a,b) ((a > b) ? (a) : (b))
#endif
#ifndef NORMALIZE
#define NORMALIZE(z) atan2(sin(z), cos(z))
#endif

namespace rbot
{
	//! OI op codes
	/*!
	 * Op codes for commands as specified by the iRobot Open Interface.
	 */
	typedef enum _OI_Opcode {

		// Command opcodes
		OI_OPCODE_START = 128,
		OI_OPCODE_BAUD = 129,
		OI_OPCODE_CONTROL = 130,
		OI_OPCODE_SAFE = 131,
		OI_OPCODE_FULL = 132,
		OI_OPCODE_POWER = 133,
		OI_OPCODE_SPOT = 134,
		OI_OPCODE_CLEAN = 135,
		OI_OPCODE_MAX = 136,
		OI_OPCODE_DRIVE = 137,
		OI_OPCODE_MOTORS = 138,
		OI_OPCODE_LEDS = 139,
		OI_OPCODE_SONG = 140,
		OI_OPCODE_PLAY = 141,
		OI_OPCODE_SENSORS = 142,
		OI_OPCODE_FORCE_DOCK = 143,
		OI_OPCODE_PWM_MOTORS = 144,
		OI_OPCODE_DRIVE_DIRECT = 145,
		OI_OPCODE_DRIVE_PWM = 146,
		OI_OPCODE_STREAM = 148,
		OI_OPCODE_QUERY = 149,
		OI_OPCODE_PAUSE_RESUME_STREAM = 150,
		OI_OPCODE_SCHEDULE_LEDS = 162,
		OI_OPCODE_DIGIT_LEDS_RAW = 163,
		OI_OPCODE_DIGIT_LEDS_ASCII = 164,
		OI_OPCODE_BUTTONS = 165,
		OI_OPCODE_SCHEDULE = 167,
		OI_OPCODE_SET_DAY_TIME = 168

	} OI_Opcode;

	/*! \class RBotInterface RBotInterface.h "include/RBotInterface.h"
	 *  \brief C++ class implementation of the RBot 100 Interface.
	 *
	 * This class implements the RBot interface protocol reverse engineered by Oleg Urzhumtcev. Based on the OpenInterface for iRobot from Goncalo Cabritas.
	 */
	class RBotInterface
	{
		public:
	
		//! Constructor
		/*!
		 * By default the constructor will set the Roomba to read only the encoder counts (for odometry).
		 *
		 *  \param new_serial_port    Name of the serial port to open.
		 *
		 *  \sa setSensorPackets()
		 */
		RBotInterface(const char * new_serial_port);
		//! Destructor
		~RBotInterface();
	
		//! Open the serial port
		/*!
		 *  \param full_init    Send full init sequence including head zeroing to RBot
		 */
		int openSerialPort(bool full_init);
		//! Close the serial port
		int closeSerialPort();
	
		//! Power down the RBot. Additional testing needed.
		int powerDown();
	
		//! Set sensor packets
		/*!
		*  Set the sensor packets one wishes to read from the roomba. By default the constructor will set the Roomba to read only the encoder counts (for odometry). 
		*
		*  \param new_sensor_packets  	Array of sensor packets to read.
		*  \param new_num_of_packets  	Number of sensor packets in the array.
		*  \param new_buffer_size		Size of the resulting sensor data reply.
		*
		*  \return 0 if ok, -1 otherwise.
		*/
		int setSensorPackets(int * new_sensor_packets, int new_num_of_packets, size_t new_buffer_size);
		//! Read sensor packets
		/*!
		*  Requested the defined sensor packets from the Roomba. If you need odometry and you requested encoder data you need to call calculateOdometry() afterwards.
		*
		*  \param timeout		Timeout in milliseconds.
		*
		* \sa calculateOdometry()
		*
		*  \return 0 if ok, -1 otherwise.
		*/
		int getSensorPackets(int timeout);
		
		//! Calculate Roomba odometry. Call after reading encoder pulses.
		void calculateOdometry();
	
		//! Drive
		/*!
		*  Send velocity commands to Roomba.
		*
		*  \param linear_speed  	Linear speed.
		*  \param angular_speed  	Angular speed.
		*
		*  \return 0 if ok, -1 otherwise.
		*/
		int drive(double linear_speed, double angular_speed);
		//! Drive direct
		/*!
		*  Send velocity commands to Roomba.
		*
		*  \param left_speed  	Left wheel speed.
		*  \param right_speed  	Right wheel speed.
		*
		*  \return 0 if ok, -1 otherwise.
		*/
		int driveDirect(int left_speed, int right_speed);

	
		//! Set the Roomba in cleaning mode. Returns the OImode to safe.
		int clean();
		//! Set the Roomba in max cleaning mode. Returns the OImode to safe.
		int max();
		//! Set the Roomba in spot cleaning mode. Returns the OImode to safe.
		int spot();
		//! Sends the Roomba to the dock. Returns the OImode to safe.
		int goDock();
	
		
	
		//! Current operation mode, one of ROOMBA_MODE_'s
		unsigned char OImode_;
	
		//! Sends the Roomba to the dock. Returns the OImode to safe.
		void resetOdometry();
		void setOdometry(double new_x, double new_y, double new_yaw);
	
		//! Roomba odometry x
		double odometry_x_;
		//! Roomba odometry y
		double odometry_y_;
		//! Roomba odometry yaw
		double odometry_yaw_;
	
		bool wall_;						//! Wall detected.
		int motor_current_[4];			//! Motor current. Indexes: LEFT RIGHT MAIN_BRUSH SIDE_BRUSH
		bool overcurrent_[4];			//! Motor overcurrent. Indexes: LEFT RIGHT MAIN_BRUSH SIDE_BRUSH
	
		unsigned char charging_state_;	//! One of OI_CHARGING_'s
		bool power_cord_;				//! Whether the Roomba is connected to the power cord or not.
		bool dock_;						//! Whether the Roomba is docked or not.
		float voltage_;					//! Battery voltage in volts.
		float current_;					//! Battery current in amps.
		char temperature_;				//! Battery temperature in C degrees.
		float charge_;					//! Battery charge in Ah.
		float capacity_;				//! Battery capacity in Ah
	
		int stasis_;					//! 1 when the robot is going forward, 0 otherwise

		private:
	
		//! Parse data
		/*!
		*  Data parsing function. Parses data comming from the Roomba.
		*
		*  \param buffer  			Data to be parsed.
		*  \param buffer_length  	Size of the data buffer.
		*
		*  \return 0 if ok, -1 otherwise.
		*/
		int parseSensorPackets(unsigned char * buffer, size_t buffer_length);
	
		int parseBumpersAndWheeldrops(unsigned char * buffer, int index);
		int parseWall(unsigned char * buffer, int index);
		int parseLeftCliff(unsigned char * buffer, int index);
		int parseFrontLeftCliff(unsigned char * buffer, int index);
		int parseFrontRightCliff(unsigned char * buffer, int index);
		int parseRightCliff(unsigned char * buffer, int index);	
		int parseVirtualWall(unsigned char * buffer, int index);
		int parseOvercurrents(unsigned char * buffer, int index);
		int parseDirtDetector(unsigned char * buffer, int index);
		int parseIrOmniChar(unsigned char * buffer, int index);
		int parseButtons(unsigned char * buffer, int index);
		int parseDistance(unsigned char * buffer, int index);
		int parseAngle(unsigned char * buffer, int index);
		int parseChargingState(unsigned char * buffer, int index);
		int parseVoltage(unsigned char * buffer, int index);
		int parseCurrent(unsigned char * buffer, int index);
		int parseTemperature(unsigned char * buffer, int index);
		int parseBatteryCharge(unsigned char * buffer, int index);
		int parseBatteryCapacity(unsigned char * buffer, int index);
		int parseWallSignal(unsigned char * buffer, int index);
		int parseLeftCliffSignal(unsigned char * buffer, int index);
		int parseFrontLeftCliffSignal(unsigned char * buffer, int index);
		int parseFontRightCliffSignal(unsigned char * buffer, int index);
		int parseRightCliffSignal(unsigned char * buffer, int index);
		int parseChargingSource(unsigned char * buffer, int index);
		int parseOiMode(unsigned char * buffer, int index);
		int parseSongNumber(unsigned char * buffer, int index);
		int parseSongPlaying(unsigned char * buffer, int index);
		int parseNumberOfStreamPackets(unsigned char * buffer, int index);
		int parseRequestedVelocity(unsigned char * buffer, int index);
		int parseRequestedRadius(unsigned char * buffer, int index);
		int parseRequestedRightVelocity(unsigned char * buffer, int index);
		int parseRequestedLeftVelocity(unsigned char * buffer, int index);
		int parseRightEncoderCounts(unsigned char * buffer, int index);
		int parseLeftEncoderCounts(unsigned char * buffer, int index);
		int parseLightBumper(unsigned char * buffer, int index);
		int parseLightBumperLeftSignal(unsigned char * buffer, int index);
		int parseLightBumperFrontLeftSignal(unsigned char * buffer, int index);
		int parseLightBumperCenterLeftSignal(unsigned char * buffer, int index);
		int parseLightBumperCenterRightSignal(unsigned char * buffer, int index);
		int parseLightBumperFrontRightSignal(unsigned char * buffer, int index);
		int parseLightBumperRightSignal(unsigned char * buffer, int index);
		int parseIrCharLeft(unsigned char * buffer, int index);
		int parseIrCharRight(unsigned char * buffer, int index);
		int parseLeftMotorCurrent(unsigned char * buffer, int index);
		int parseRightMotorCurrent(unsigned char * buffer, int index);
		int parseMainBrushMotorCurrent(unsigned char * buffer, int index);
		int parseSideBrushMotorCurrent(unsigned char * buffer, int index);
		int parseStasis(unsigned char * buffer, int index);
	
		//! Buffer to signed int
		/*!
		*  Parsing helper function. Converts 2 bytes of data into a signed int value. 
		*
		*  \param buffer  	Data buffer.
		*  \param index  	Position in the buffer where the value is.
		*
		*  \sa buffer2unsigned_int()
		*
		*  \return A signed int value.
		*/
		int buffer2signed_int(unsigned char * buffer, int index);
		//! Buffer to unsigned int
		/*!
		*  Parsing helper function. Converts 2 bytes of data into an unsigned int value. 
		*
		*  \param buffer  	Data buffer.
		*  \param index  	Position in the buffer where the value is.
		*
		*  \sa buffer2signed_int()
		*
		*  \return An unsigned int value.
		*/
		int buffer2unsigned_int(unsigned char * buffer, int index);
	
		//! Start OI
		/*!
		*  Start the OI, change to roomba to a OImode that allows control.
		*
		*  \param full_control    Whether to set the Roomba on OImode full or not.
		*
		*  \return 0 if ok, -1 otherwise.
		*/
		int startOI(bool full_control);
		//! Send OP code
		/*!
		*  Send an OP code to Roomba.
		*
		*  \param code  			OP code to send.
		*
		*  \return 0 if ok, -1 otherwise.
		*/
		int sendOpcode(OI_Opcode code);
	
		//! Serial port to which the robot is connected
		std::string port_name_;
		//! Cereal port object
		cereal::CerealPort * serial_port_;
	
		//! Stream variable. NOT TESTED
		bool stream_defined_;

		//! Amount of distance travelled since last reading. Not being used, poor resolution. 
		int distance_;
		//! Amount of angle turned since last reading. Not being used, poor resolution. 
		int angle_;
		//! Delta encoder counts.
		int encoder_counts_[2];
		//! Last encoder counts reading. For odometry calculation.
		uint16_t last_encoder_counts_[2];
	};

}

// EOF
