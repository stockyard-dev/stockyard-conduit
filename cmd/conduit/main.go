package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-conduit/internal/server";"github.com/stockyard-dev/stockyard-conduit/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9690"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./conduit-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("conduit: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Conduit — data sync engine\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("conduit: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
