package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Pipe struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Source string `json:"source"`
	Destination string `json:"destination"`
	Transform string `json:"transform"`
	Enabled string `json:"enabled"`
	ThroughputPerSec int `json:"throughput"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"conduit.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS pipes(id TEXT PRIMARY KEY,name TEXT NOT NULL,source TEXT DEFAULT '',destination TEXT DEFAULT '',transform TEXT DEFAULT '',enabled TEXT DEFAULT 'true',throughput INTEGER DEFAULT 0,created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Pipe)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO pipes(id,name,source,destination,transform,enabled,throughput,created_at)VALUES(?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Source,e.Destination,e.Transform,e.Enabled,e.ThroughputPerSec,e.CreatedAt);return err}
func(d *DB)Get(id string)*Pipe{var e Pipe;if d.db.QueryRow(`SELECT id,name,source,destination,transform,enabled,throughput,created_at FROM pipes WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Source,&e.Destination,&e.Transform,&e.Enabled,&e.ThroughputPerSec,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Pipe{rows,_:=d.db.Query(`SELECT id,name,source,destination,transform,enabled,throughput,created_at FROM pipes ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Pipe;for rows.Next(){var e Pipe;rows.Scan(&e.ID,&e.Name,&e.Source,&e.Destination,&e.Transform,&e.Enabled,&e.ThroughputPerSec,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM pipes WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM pipes`).Scan(&n);return n}
