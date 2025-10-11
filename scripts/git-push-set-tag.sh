# 说明: Git提交当前分支并设置tag
# 运行方式: sh git-push-set-tag.sh 提交代码的备注信息 tag名称

if [ ! -n "$1" ] ;then
    # mark="修改"
    echo "参数1:备注不能为空"
    exit 2
else
    mark=$1
fi

if [ ! -n "$2" ] ;then
    # 获取最新的 tag
    latest_tag=$(git tag --sort=-creatordate | head -n 1)
    if [ -n "$latest_tag" ]; then
        echo "参数2:tag名称不能为空,当前最新tag: $latest_tag"
    else
        echo "参数2:tag名称不能为空,当前没有可用的tag,你可以使用类似 v0.0.1 的格式作为第一个tag"
    fi
    exit 2
else
    tag=$2
fi

git pull
git add .
git commit -m "$mark"
git push

git tag -a "$tag" -m "new tag $tag"
git push origin "$tag"
