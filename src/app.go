import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Client
}

func (a *App) Initialize(user, password, dbName string) {}
func (a *App) Run(addr string)                          {}