#ifndef __EDR_COMMUNICATION_H__
#define __EDR_COMMUNICATION_H__

//指定了数据类型
enum {
	EDR_ATTR_UNSPEC,
	EDR_ATTR_TYPE,		/* string */
	EDR_ATTR_TIMESTAMP,	/* u64 */

	EDR_ATTR_CUR_PID,	/* u64 */
	EDR_ATTR_CHILD_PID,	/* u64 */
	EDR_ATTR_TASKNAME,	/* string */
	EDR_ATTR_FILENAME,	/* string */
	__EDR_ATTR_MAX
};
#define EDR_ATTR_MAX (__EDR_ATTR_MAX - 1)
#define EDR_FAMILY_NAME "EdrFamily"
#define EDR_FAMILY_VERSION 0x1


//命令类型
enum {
    TEST_CMD_UNSPEC,
    TEST_CMD_MSG,
    TEST_CMD_MSG2,
    __TEST_CMD_MAX,
};
#define TEST_CMD_MAX (__TEST_CMD_MAX - 1)




#endif	//__EDR_COMMUNICATION_H__
