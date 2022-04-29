#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>

#include <net/genetlink.h>
#include "edr_communication.h"
static int msg_func(struct sk_buff* skb2, struct genl_info* info);
static int msg2_func(struct sk_buff* skb, struct genl_info* info);


//定义ATTR类型对应的数据类型
static const struct nla_policy edr_nl_policy[EDR_ATTR_MAX + 1] = {
	[EDR_ATTR_TYPE] = { .type = NLA_STRING },
	[EDR_ATTR_TIMESTAMP] = { .type = NLA_U64},

	[EDR_ATTR_CUR_PID] = { .type = NLA_U64},
	[EDR_ATTR_CHILD_PID] = { .type = NLA_U64},

	[EDR_ATTR_TASKNAME] = { .type = NLA_STRING},
	[EDR_ATTR_FILENAME] = { .type = NLA_STRING},
};
//定义EDR CMD
static struct genl_ops edr_ops[] = {
	{
		.cmd = TEST_CMD_MSG,
		.flags = 0,
		.doit = msg_func,
		.policy = edr_nl_policy,
		.dumpit = NULL,
	},
	{
		.cmd = TEST_CMD_MSG2,
		.flags = 0,
		.doit = msg2_func,
		.policy = edr_nl_policy,
		.dumpit = NULL,
	},
};

static struct genl_family edr_nl_family __ro_after_init = {
	.module = THIS_MODULE,
	.name = EDR_FAMILY_NAME,
	.version = EDR_FAMILY_VERSION,
	.maxattr = EDR_ATTR_MAX,
	.policy = edr_nl_policy,

	.ops = edr_ops,
	.n_ops = ARRAY_SIZE(edr_ops),

	.hdrsize = 0,
};

static int msg_func(struct sk_buff* skb2, struct genl_info* info)
{
	printk("msg_func ...\n");
	struct nlattr *na;
    struct sk_buff *skb;
	int rc;
	void *msg_hdr;
	char *data;
	if(info == NULL)
			goto error;
	//对于每个属性，genl_info的域attrs可以索引到具体结构，里面有payload
	na = info->attrs[EDR_ATTR_FILENAME];
	if(na){
		data = (char *)nla_data(na);
		if(!data) printk("Receive data error!\n");
		else printk("Recv:%s\n",data);
	}else{
		printk("No info->attrs %s\n","EDR_ATTR_FILENAME");
	}

	skb = genlmsg_new(NLMSG_GOODSIZE,GFP_KERNEL);
	if(!skb) goto error;

	/*构建消息头，函数原型是
	genlmsgput(struct sk_buff *,int pid,int seq_number,
			struct genl_family *,int flags,u8 command_index);
	*/
	msg_hdr = genlmsg_put(skb,0,info->snd_seq+1, &edr_nl_family,
							0,TEST_CMD_MSG);
	if(msg_hdr == NULL){
		rc = -ENOMEM;
		goto error;
	}

	//填充具体的netlink attribute:DOC_EXMPL_A_MSG，这是实际要传的数据
	rc = nla_put_string(skb,EDR_ATTR_FILENAME,"HelloWorld");
	if(rc != 0) goto error;

	genlmsg_end(skb,msg_hdr);//消息构建完成
	//单播发送给用户空间的某个进程
	//rc = genlmsg_unicast(genl_info_net(info),skb,info->snd_portid);
	printk("1 : %x\n", (long)genl_info_net(info));
	printk("2 : %x\n", (long)&init_net);
	rc = genlmsg_unicast(&init_net,skb,info->snd_portid);
	if(rc != 0){
			printk("Unicast to process:%d failed!\n",info->snd_portid);
			goto error;
	}
	return 0;

error:
	printk("Error occured in doc_echo!\n");
	return 0;
}
static int msg2_func(struct sk_buff* skb, struct genl_info* info)
{
	printk("msg2_func ...\n");
	return 0;
}


static int __init edr_communication_init(void)
{
	int ret;
	ret = genl_register_family(&edr_nl_family);
	if (ret) {
		printk("genl_register_family failed\n");
		return -1;
	}
	printk("genl_register_family success\n");
	return 0;
}



static void __exit edr_communication_exit(void)
{
	printk("genl_unregister_family...\n");
	genl_unregister_family(&edr_nl_family);
}


module_init(edr_communication_init);
module_exit(edr_communication_exit);
MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Monitor Syscall sys_execve");
