//
// hash_test.c
// taka-wang
//

#include <stdio.h>
#include <stdlib.h>

#include <unistd.h>
#include <string.h>
#include <time.h>
#include <errno.h>


#include "../../json.h"

int enable_syslog = 1;

//=================================
// Test functions
//=================================


// decode json string
void test_json_decode()
{
    char jstr[] = "{\n"
    "    \"ip\": \"192.168.3.2\",\n"
    "    \"port\": \"502\",\n"
    "    \"slave\": 22,\n"
    "    \"tid\": 1,\n"
    "    \"mode\": \"tcp\",\n"
    "    \"cmd\": \"fc5\",\n"
    "    \"addr\": 250,\n"
    "    \"len\": 10,\n"
    "    \"data\": [1,2,3,4]\n"
    "}";
    cJSON *json = cJSON_Parse(jstr);
    if (json)
    {
        printf("%s\n", cJSON_Print(json));
        printf("---------\n");
        printf("ip:%s\n", cJSON_GetObjectItem(json, "ip")->valuestring);
        printf("port:%s\n",cJSON_GetObjectItem(json, "port")->valuestring);
        printf("mode:%s\n", cJSON_GetObjectItem(json, "mode")->valuestring);
        printf("addr:%d\n",cJSON_GetObjectItem(json, "addr")->valueint);
        
        // handle array
        cJSON * data = cJSON_GetObjectItem(json, "data");
        for (int i = 0 ; i < cJSON_GetArraySize(data) ; i++)
        {
            int subitem = cJSON_GetArrayItem(data, i)->valueint;
            printf("idx:%d,v:%d\n", i, subitem);
        }
        if (json != NULL) cJSON_Delete(json);
    }
}

// encode to json string
void test_json_encode()
{
    int mdata[4]={116,943,234,38793};
    cJSON *root;
    root = cJSON_CreateObject();
    cJSON_AddNumberToObject(root, "tid", 14672035611234);
    cJSON_AddItemToObject(root,"data", cJSON_CreateIntArray(mdata,4));
    cJSON_AddStringToObject(root, "status", "ok");
    printf("%s\n", cJSON_Print(root));
    cJSON_Delete(root);
    /*
    {
        "tid": 22,
        "data": [1,2,3,4],
        "status": "ok"
    }
    */
}



// ENTRY
int main()
{


    test_json_decode();
    test_json_encode();


}