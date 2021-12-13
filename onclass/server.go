package onclass

import "net/http"

//
//type Server interface {
//	Route(pattern string, handlerFunc http.HandlerFunc)
//	Start(address string) error
//}
//
//type sdkHttpServer struct {
//	Name string
//}
//
//// 结构体作为接口的方法接收器，最好都是用指针的形式
//func (s *sdkHttpServer) Route(pattern string, handlerFunc http.HandlerFunc) {
//	http.HandleFunc(pattern, handlerFunc)
//}
//
//func (s *sdkHttpServer) Start(address string) error {
//	return http.ListenAndServe(address, nil)
//}
//
//func NewHttpServer(name string) Server {
//	return &sdkHttpServer{	//当返回实际类型所实现的接口的时候，需要返回指针
//		Name:name,
//	}
//}
//
////若是type A B形式，就不使用指针形式，但这种方法一般不使用
//type Handle func()
//func (h Handle) Hello(){
//
//}
//
//func SignUpWithoutContext(w http.ResponseWriter, r *http.Request) {
//	req := &signUpReq{}
//
//	////////////////// 么有context时，使用原生的方法读json文件等处理 ////////////////
//	//body, err := io.ReadAll(r.Body)
//	//if err != nil {
//	//	fmt.Fprintf(w, "read body failed: %v", err)
//	//	// 要返回掉，不然就会继续执行后面的代码
//	//	return
//	//}
//	//err = json.Unmarshal(body, req)
//	//if err != nil {
//	//	fmt.Fprintf(w, "deserialized failed: %v", err)
//	//	// 要返回掉，不然就会继续执行后面的代码
//	//	return
//	//}
//	/////////////////////////////////////////////////////////////////////////
//
//	// 使用Context 处理json
//	ctx := Context{
//		W: w,
//		R: r,
//	}
//	err := ctx.ReadJson(req)
//	if err != nil{
//		fmt.Fprintf(w,"error %v", err)
//	}
//
//
//	resp := commonResponse{
//		Data: 123,
//	}
//
//	//////////////////// 没有Context时response的处理 //////////////////
//	//rep, err := json.Marshal(resp)
//	//if err != nil{
//	//	fmt.Fprintf(w, "error %v ", err)
//	//}
//	////////////////////////////////////////////////////////////////
//
//	err = ctx.OkJson(resp)
//	if err != nil{
//		fmt.Fprintf(w, "error %v ", err)
//	}
//
//	// 返回一个虚拟的 user id 表示注册成功了
//	// fmt.Fprintf(w, string(rep))	//Marshal返回的是byte的切片，而这里使用Fprintf输出的是string类型，所以要类型转换
//}

type Server interface {
	Route(method string, pattern string, handleFunc func(ctx *Context)) //添加method
	Start(address string) error
}

type sdkHttpServer struct {
	Name    string
	handler *HandlerBasedMap //添加HandlerBasedMap
}

func NewHttpServer(name string) Server {
	return &sdkHttpServer{ //当返回实际类型所实现的接口的时候，需要返回指针
		Name: name,
	}
}

//func (s *sdkHttpServer) Route(pattern string,
//	handleFunc func(ctx *Context)) {
//	http.HandleFunc(pattern, func(write http.ResponseWriter, request *http.Request){
//		ctx := NewContext(write, request)
//		handleFunc(ctx)
//	})
//}
//func (s *sdkHttpServer) Start(address string) error {
//	return http.ListenAndServe(address, nil)
//}

///////////////////// 添加Method /////////////////////
func (s *sdkHttpServer) Route(
	method string,
	pattern string,
	handleFunc func(ctx *Context)) {
	key := s.handler.Key(method, pattern)
	s.handler.handlers[key] = handleFunc
}
func (s *sdkHttpServer) Start(address string) error {
	http.Handle("/", s.handler)
	return http.ListenAndServe(address, nil)
}

/////////////////////////////////////////////////////

func SignUp(ctx *Context) {
	resp := &signUpReq{}
	err := ctx.ReadJson(resp)
	if err != nil {
		err := ctx.BadRequestJson(err)
		if err != nil {
			return
		}
	}

	rep := commonResponse{
		Data: 123,
	}
	err = ctx.OkJson(rep)
	if err != nil {
		err := ctx.BadRequestJson(err)
		if err != nil {
			return
		}
	}
}

type signUpReq struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}
