function htmlEncode(str){
    if (typeof(str)=="undefined" || str == null ){
        return "";
    }

    if(typeof(str) != "string"){
        
        str = str.toString();

    }
     /*str = str.replace(/&/g,'&amp;').replace(/</g,'&lt;',/>/g,'&gt;').replace(/"/g,'&quot;').replace(/'/g,'&#39;');
     return str*/
    return str.replace(/&/g,'&amp;')
            .replace(/</g,'&lt;',/>/g,'&gt;')
            .replace(/"/g,'&quot;')
            .replace(/'/g,'&#39;');
 
}

function ajaxRequest(method,url,params,xsrfToken,callback){
   //console.log("2. ajaxRequest")
        jQuery.ajax({
            type: method,
            url: url,
            data: params,
            
            //xhr为ajax的原生对象，在发送请求之前能够在请求头部设置新的xsrf的token
            beforeSend: function(xhr){ //经测试，无法获取到cookie中的_xsrf的token，还未找到原因     
                //xhr.setRequestHeader("X-Xsrftoken", jQuery.base64.atob(jQuery.cookie("_xsrf").split("|")[0]) );
                if(xsrfToken !=""){  //注意: 新建和编辑的xsrf的token在提交的html表单中因此无需传入xsrfToken ,而锁定\删除控制按钮等设置了从dailog-body中获取xsrf,并赋值给xsrfToken
                    xhr.setRequestHeader("X-Xsrftoken",xsrfToken);
                }       
      
            },
            success: function(response){
                
                switch(response["code"]){
                    case 200:
                    // console.log("3. ajaxRequest success")
                        //callback(response);//只有调用此callback函数，dailog-js.html中的引入的function(response)才能生效
                        //alert("成功");

                        swal("操作成功！",response["text"],"success")
                        jQuery.notify({    
                            title: 'CMDB管理系统：',
                            icon:"/static/image/title.jpg",
                            message: response["text"]
                        },
                        {
                            type: "success",
                            icon_type:"image"
                        
                        })
                        
                        // UserPageTables.ajax.reload(null,false);//ajax自带的刷新，ajax.reload( callback, resetPaging )：null表示不回调函数，false表示刷新不重置到第一页
                        callback(response);//只有调用此callback函数，dailog-js.html中的引入的function(response)才能生效

                        break;//中断继续
                    
                    case 304:
                        swal("提示！",response.result[0].Message,"info")
                        break;
                    case 400:
                        var errors = []
                        jQuery.each(response["result"],function(index,value){
                            errors.push(value["Message"])
                        });
                       // console.log("hahah",errors)
                        if(errors.length == 0){
                            errors.push(response["text"])  
                        }
                        swal("出错了！",errors,"error");

                        //alert(errors.join("\n"));
                        break;
                    case 403:
                        swal("操作禁止！",response["text"],"warning")
                        break;
                    case 500:
                        alert("服务器崩溃");
                        break;
                    default:
                        alert("未知错误");
                        break;
                }
            },
            dataType: "json"
        })

    

 
}