#ARCH
OS=`uname -s`
currPWD=$(cd `dirname $0`; pwd)

echo "currPWD [$currPWD] OS [$OS]"


echo "config go env"
export GOPATH=$GOPATH:$currPWD
export GO111MODULE=on
export GOPROXY=https://goproxy.io

#install module
echo "go mod fly3d engine"
cd $currPWD/fly3d
go get

#build shader
shaderfile=$currPWD/fly3d/module/effects/effectstore.go
echo "" > $shaderfile
echo "package effects\n" >> $shaderfile
echo "var ShadersStore = map[string]string{}\n" >> $shaderfile
echo "func init() { \n">> $shaderfile
for element in `ls $currPWD/fly3d/shaders`
    do  

        file_name=`echo $element|awk -F"." '{print $1}'`  
        type_name=`echo $element|awk -F"." '{print $2}'`  
       
        data_content=`cat $currPWD/fly3d/shaders/$element`

        echo "ShadersStore[\"${file_name}_${type_name}\"] = \`$data_content\` \n" >>  $shaderfile
  


    done

echo "} \n" >> $shaderfile