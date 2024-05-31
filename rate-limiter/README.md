# Objetivo

Limitar o tráfego em um servidor Web em Go.

# Funcionamento
 
A biblioteca fornece um interceptador de chamadas do servidor Web que adiciona a contagem de requisições e seu bloqueio quando a quantidade ultapassar um limite.

# Configuração

1. Configuração no arquivo .env

- período de contagem de requisições:

    PERIOD=1

- limite de requisições por período sem o uso de token:

    LIMIT=10
 
2. Configuração em NewLimitedHandler()

- duração do bloqueio após ultrapassar o limite de requisições:

	blockInSeconds int

- nome dos tokens:

	tokens      [2]string

- limite de requisições dos tokens:

	tokenLimits [2]int

# Uso

Importar o package:

    import "github.com/gsusin/goexpert/rate-limiter/limiter"

Criar um tipo contendo uma referência a LimitedHandler:

    type LimitedServer struct {
	    lim *limiter.LimitedHandler
    }
    
Criar um método de interceptação que direciona para LimitedFunc() e fornece o handler final:

    func (ls *LimitedServer) LimitedGetCourseId(w http.ResponseWriter, r *http.Request) {
	    ls.lim.LimitedFunc(w, r, GetCourseId)(w, r)
    }
 
Usar NewLimitedServer para criar o handler:

    func init() {
	    mux := http.NewServeMux()
	    ls := LimitedServer{}
	    ls.lim = limiter.NewLimitedHandler()
	    mux.HandleFunc("/course", ls.LimitedGetCourseId)
	    go http.ListenAndServe(":8080", mux)
    }
