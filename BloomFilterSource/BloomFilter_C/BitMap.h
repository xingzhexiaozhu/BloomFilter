#pragma once  
  
#include<stdint.h>  
#include<stddef.h>  
  
typedef struct BitMap{  
    uint64_t* data;  
    size_t capacity;  
}BitMap;  
  
/*初始化*/  
void BitMapInit(BitMap* bm, size_t capacity);  
  
/*把index位设置为１*/  
void BitMapSet(BitMap* bm, size_t index);  
  
/*把index位设置为０*/  
void BitMapUnSet(BitMap* bm, size_t index);  
  
/*测试index为１还是为０，如果是１，就返回１，否则返回０*/  
int BitMapTest(BitMap* bm, size_t index);  
  
/*把这个位图所有位都设置为１*/  
void BitmapFill(BitMap* bm);  
  
/*把整个位图所有位都设置为０*/  
void BitMapClear(BitMap* bm);  
  
/*销毁*/  
void BitMapDestory(BitMap* bm);  
