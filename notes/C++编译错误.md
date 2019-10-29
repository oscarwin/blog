**问题：**linux 生成动态库时提示relocation R_X86_64_32 against `__gxx_personality_v0' can not be used when making a shared object;

图片：
![](https://github.com/oscarwin/blog/blob/master/image/compile_error.png)

**解决：**在编译add.o文件和sub.o时也应该加上-fPIC选项。原因：因为-fPIC选项让动态链接库生成地址无关的代码，动态链接库中所有函数地址应该都有一定的规则来选择，因此在编译动态库的依赖文件时也要加上该选项。（个人理解）

---

**问题：**今天使用sort排序一个结构体，编译一直不过。

最开始是将自定义的比较函数放在头文件中
**解决：**将比较函数放入cpp中，且声明为const类型的参数传入
```
bool _MyComparisonFun(const boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo& a, const boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo& b)
{
    return a.m_uiStatus < b.m_uiStatus;
}
```
```
/usr/local/lib/gcc/x86_64-unknown-linux-gnu/4.4.7/../../../../include/c++/4.4.7/bits/stl_algo.h: In function ‘const _Tp& std::__median(const _Tp&, const _Tp&, const _Tp&, _Compare) [with _Tp = boss::shopgroup::shopgroup_shop_struct::_SHOPGROUP_ACTIVE_INFO, _Compare = bool (*)(boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo&, boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo&)]’:
/usr/local/lib/gcc/x86_64-unknown-linux-gnu/4.4.7/../../../../include/c++/4.4.7/bits/stl_algo.h:2301:   instantiated from ‘void std::__introsort_loop(_RandomAccessIterator, _RandomAccessIterator, _Size, _Compare) [with _RandomAccessIterator = __gnu_cxx::__normal_iterator<boss::shopgroup::shopgroup_shop_struct::_SHOPGROUP_ACTIVE_INFO*, std::vector<boss::shopgroup::shopgroup_shop_struct::_SHOPGROUP_ACTIVE_INFO, std::allocator<boss::shopgroup::shopgroup_shop_struct::_SHOPGROUP_ACTIVE_INFO> > >, _Size = long int, _Compare = bool (*)(boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo&, boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo&)]’
/usr/local/lib/gcc/x86_64-unknown-linux-gnu/4.4.7/../../../../include/c++/4.4.7/bits/stl_algo.h:5258:   instantiated from ‘void std::sort(_RAIter, _RAIter, _Compare) [with _RAIter = __gnu_cxx::__normal_iterator<boss::shopgroup::shopgroup_shop_struct::_SHOPGROUP_ACTIVE_INFO*, std::vector<boss::shopgroup::shopgroup_shop_struct::_SHOPGROUP_ACTIVE_INFO, std::allocator<boss::shopgroup::shopgroup_shop_struct::_SHOPGROUP_ACTIVE_INFO> > >, _Compare = bool (*)(boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo&, boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo&)]’
../../../storage/mysql/shop/dao_shopgroup_shop.cpp:217:   instantiated from here
/usr/local/lib/gcc/x86_64-unknown-linux-gnu/4.4.7/../../../../include/c++/4.4.7/bits/stl_algo.h:124: error: invalid initialization of reference of type ‘boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo&’ from expression of type ‘const boss::shopgroup::shopgroup_shop_struct::_SHOPGROUP_ACTIVE_INFO’
/usr/local/lib/gcc/x86_64-unknown-linux-gnu/4.4.7/../../../../include/c++/4.4.7/bits/stl_algo.h:125: error: invalid initialization of reference of type ‘boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo&’ from expression of type ‘const boss::shopgroup::shopgroup_shop_struct::_SHOPGROUP_ACTIVE_INFO’
/usr/local/lib/gcc/x86_64-unknown-linux-gnu/4.4.7/../../../../include/c++/4.4.7/bits/stl_algo.h:127: error: invalid initialization of reference of type ‘boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo&’ from expression of type ‘const boss::shopgroup::shopgroup_shop_struct::_SHOPGROUP_ACTIVE_INFO’
/usr/local/lib/gcc/x86_64-unknown-linux-gnu/4.4.7/../../../../include/c++/4.4.7/bits/stl_algo.h:131: error: invalid initialization of reference of type ‘boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo&’ from expression of type ‘const boss::shopgroup::shopgroup_shop_struct::_SHOPGROUP_ACTIVE_INFO’
/usr/local/lib/gcc/x86_64-unknown-linux-gnu/4.4.7/../../../../include/c++/4.4.7/bits/stl_algo.h:133: error: invalid initialization of reference of type ‘boss::shopgroup::shopgroup_shop_struct::stShopgroupActiveInfo&’ from expression of type ‘const boss::shopgroup::shopgroup_shop_struct::_SHOPGROUP_ACTIVE_INFO’
```

---

**问题：**使用set时，当调用insert成员函数时就会出错

**解决：**具体的出错信息忘记保存了，出错的原因在于set中的second对象是一个自定义对象，而set的底层实现时红黑树，红黑树插入时需要比较大小，而自定义的对象没有重载<符号，因此到此编译错误。对于关联容器map和set，key如果是自定义对象，则都需要重载<符号。为什么只需要重载小于符号？因为其他的符号都可以通过小于符号实现。

```
bool opeartor < (const T& obj)
{
    return this->item < obj->item;
}
```

---

**问题：**discards qualifiers

**解决：**discards qualifiers错误是由于const成员调用了非const函数。

```
int CCutPriceTemplateMsgHandler::SendWeiXinMsg(const CPushMsgqStruct& stPushMsg)
{
    boss::serialize::Buffer oBuf;
    oBuf<<stPushMsg;
    string sMsg;
    sMsg.assign(oBuf.getbuf(), oBuf.length());
    uint32_t uiRetSeq = 0;
    int iRet = m_oProducer.MultiSendMsg(sMsg, uiRetSeq);
    if(iRet < 0)
    {    
        C2C_WW_LOG_ERR(iRet, "MultiSendMsg failed, iRet[%d], errMsg[%s]", iRet, m_oProducer.GetLastErrMsg());
        UMP->PerStop(strPerKey, 1);
        return -1;
    }

    if(stPushMsg.ulMsgType == WEIXIN_CUTPRICE_JOIN)
    {
        C2C_WW_LOG("Send cutprice half msg succ to JD_WEIXIN, strPin[%s], ddwActiveId[%s], strFirst[%s], strKey1[%s], strKey2[%s], strKey3[%s], strKey4[%s],           
        remark[%s]",stPushMsg.sPin.c_str(), stPushMsg.mapOtherData["activeId"].c_str(), stPushMsg.mapOtherData["first"].c_str(),  
        stPushMsg.mapOtherData["keyword1"].c_str(),stPushMsg.mapOtherData["keyword2"].c_str(), stPushMsg.mapOtherData["keyword3"].c_str(), 
        stPushMsg.mapOtherData["keyword4"].c_str(),stPushMsg.mapOtherData["remark"].c_str());
    }

    return 0;
}
```

在上面的例子中传入了const的引用对象，但是map使用了操作符[]访问数据，而[]没有const版本，因此报错。对于C++11而言可以通过at函数来解决，at是const函数。

错误信息如下：
```
cutprice_template_msg.cpp:670: error: passing ‘const std::map<std::basic_string<char, std::char_traits<char>, std::allocator<char> >, std::basic_string<char, std::char_traits<char>, std::allocator<char> >, std::less<std::basic_string<char, std::char_traits<char>, std::allocator<char> > >, std::allocator<std::pair<const std::basic_string<char, std::char_traits<char>, std::allocator<char> >, std::basic_string<char, std::char_traits<char>, std::allocator<char> > > > >’ as ‘this’ argument of ‘_Tp& std::map<_Key, _Tp, _Compare, _Alloc>::operator[](const _Key&) [with _Key = std::basic_string<char, std::char_traits<char>, std::allocator<char> >, _Tp = std::basic_string<char, std::char_traits<char>, std::allocator<char> >, _Compare = std::less<std::basic_string<char, std::char_traits<char>, std::allocator<char> > >, _Alloc = std::allocator<std::pair<const std::basic_string<char, std::char_traits<char>, std::allocator<char> >, std::basic_string<char, std::char_traits<char>, std::allocator<char> > > >]’ discards qualifiers
```

