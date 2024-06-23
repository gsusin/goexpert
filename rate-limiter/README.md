# Objetivo

Limitar o tráfego em um servidor Web em Go.

# Funcionamento
 
A biblioteca fornece um interceptador de chamadas do servidor Web que adiciona a contagem de requisições e seu bloqueio quando a quantidade ultapassar uma taxa limite. Podem ser usados tokens para habilitar taxas limites diferentes da padrão.

# Configuração

1. Configuração no arquivo .env

- período de contagem de requisições, em segundos:

    PERIOD=1

- limite de requisições por período sem o uso de token:

    LIMIT=10
 
- duração do bloqueio após ultrapassar o limite de requisições, em segundos:

    BLOCK=1

- limite de requisições por período de tempo para o token "HIGH":

    LIMIT_TOKEN_HIGH=20
    
- limite de requisições por período de tempo para o token "LOW":

    LIMIT_TOKEN_LOW=20
    
2. Configuração em NewLimitedHandler()

- nome dos tokens:

	tokens      [2]string

# Uso

1. Importar o package:

    import "github.com/gsusin/goexpert/rate-limiter/limiter"

2. Criar um tipo contendo uma referência a LimitedHandler:

    type LimitedServer struct {
	    lim *limiter.LimitedHandler
    }
    
3. Criar um método de interceptação que direciona para LimitedFunc() e fornece o handler final:

    func (ls *LimitedServer) LimitedGetCourseId(w http.ResponseWriter, r *http.Request) {
	    ls.lim.LimitedFunc(w, r, GetCourseId)(w, r)
    }
 
4. Usar NewLimitedServer para criar o handler, passando o mecanismo de storage RedisStorage:

    func init() {
	    mux := http.NewServeMux()
	    ls := LimitedServer{}
	    as := limiter.NewRedisStorage()
	    ls.lim = limiter.NewLimitedHandler(&as)
	    mux.HandleFunc("/course", ls.LimitedGetCourseId)
	    go http.ListenAndServe(":8080", mux)
    }

Para trocar o mecanismo de storage para MemoryStorage, usar:

    as := limiter.NewMemoryStorage()
        
# Testes manuais

1. Configurar o arquivo .env com limites adequados aos testes manuais.

2. Subir o Redis:

    docker-compose up -d

3. Executar o servidor de testes:

    go run cmd/main.go

4. Acessar uma das URLs com frequências variáveis:

    http://localhost:8080/course
    http://localhost:8080/course?code=LOW
    http://localhost:8080/course?code=HIGH
