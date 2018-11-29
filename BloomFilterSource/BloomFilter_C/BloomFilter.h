#pragma once  
  
#include<stddef.h>  
#include"BitMap.h"  
  
#define HASHFUNCMAXSIZE 2  
#define BitMapCapacity 1024  
  
  
/*字符串哈希函数*/  
typedef size_t(*HashFunc)(const char*);  
  
typedef struct BloomFilter{  
    BitMap bitmap;  
    HashFunc hash_func[HASHFUNCMAXSIZE];  
}BloomFilter;  
  
  
  
/*初始化布隆过滤器*/  
void BloomFilterInit(BloomFilter* bf);  
  
/*插入*/  
void BloomFilterInsert(BloomFilter* bf, const char* str);  
  
int BloomFilterExit(BloomFilter* bf, const char* str);  
  
void BloomFilterDestory(BloomFilter* bf);  
  
/*按照当前的设计是不允许删除的*/  
