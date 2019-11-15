# Filter
基于bloom filter实现BloomFilterTool工具，该工具功能：      
1、实现业界的去重功能；                       
2、版本切换功能（可按时间切换、按阈值切换）；          
3、统计各版本下bloomfilter中已有元素个数；           
用法参考filter_test.go         

# Bloom Filter
总结Bloom Filter的go实现版本、Java实现版本和C实现版本            
同时将Bloom Filter的优化版本Cuckoo Filter也总结进来


# 参考的实现：
Bloom Filter(go):   
   	https://godoc.org/github.com/willf/bloom   
   	https://github.com/willf/bloom.git         
Cuckoo Filter(go):   
   	https://github.com/irfansharif/cfilter.git      
Bloom Filter(Java):   
   	https://github.com/JueFan/BloomFilterTool.git     
Bloom Filter(C):   
   	https://blog.csdn.net/Qregi/article/details/79443499
