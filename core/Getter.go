package Rigo

//我超，这里差点让我大脑宕机了
//想了一会，原来是用一个自定义的函数作为变量去实现了一个接口
//然后这个自定义函数变量调用Get方法的时候可以调用自己
//什么鬼，自己调用自己，但仔细想一下，常见http框架里面的HandlerFunc也是这样的吧.
//查了一下，这就是[接口型函数](https://geektutu.com/post/7days-golang-q1.html)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
