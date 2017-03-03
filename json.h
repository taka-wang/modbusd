/**
 * @file json.h
 * @author Taka Wang
 * @brief cJSON helper functions header
*/

#pragma once

#include "json/cJSON.h"

/**
 * @brief Get char string via key from cJSON object
 *
 * @param inJson cJSON object.
 * @param key Json key.
 * @return C char string.
 */
char * json_get_char (cJSON *inJson, const char *key);

/**
 * @brief Get integer value via key from cJSON object
 *
 * @param inJson cJSON object.
 * @param key Json key.
 * @return Integer.
 */
int json_get_int (cJSON *inJson, const char *key);


/**
 * @brief Set integer value via key to existed cJSON object
 *
 * @param inJson cJSON object.
 * @param key Json key.
 * @param value Integer value.
 * @return Void.
 */
void json_set_int (cJSON *inJson, const char *key, int value);

/**
 * @brief Get double integer value via key from cJSON object
 *
 * @param inJson cJSON object.
 * @param key Json key.
 * @return Double integer value.
 */
double json_get_double (cJSON *inJson, const char *key);

/**
 * @brief Set double integer value via key to existed cJSON object
 *
 * @param inJson cJSON object.
 * @param key Json key.
 * @param double Double integer value.
 * @return Void.
 */
void json_set_double (cJSON *inJson, const char *key, double value);

/**
 * @brief Get long integer value via key from cJSON object
 *
 * @param inJson cJSON object.
 * @param key Json key.
 * @return Long.
 */
long json_get_long (cJSON *inJson, const char *key);

/**
 * @brief Load JSON file to cJSON object
 *
 * @param fname File name string.
 * @param outJson Pointer to cJSON output object.
 * @return Success or not. 
 */ 
int file_to_json (const char *fname, cJSON **outJson);

/**
 * @brief Save cJSON object to JSON file
 *
 * @param fname File name string.
 * @param inJson cJSON input object.
 * @return Success or not. 
 */
int json_to_file (const char *fname, cJSON *inJson);
