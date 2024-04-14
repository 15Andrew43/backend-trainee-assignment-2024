import os
import sys
from dotenv import load_dotenv


load_dotenv()


num_inserts = int(sys.argv[1])


features_inserts = "INSERT INTO features (name) VALUES "
for i in range(1, num_inserts + 1):
    features_inserts += f"('feature_{i}'),\n"


tags_inserts = "INSERT INTO tags (name) VALUES "
for i in range(1, num_inserts + 1):
    tags_inserts += f"('tag_{i}'),\n"


postgres_host = "localhost"
postgres_port = os.getenv("POSTGRES_PORT")
postgres_user = os.getenv("POSTGRES_USER")
postgres_db = os.getenv("POSTGRES_DB")
postgres_password = os.getenv("POSTGRES_PASSWORD")


if not all([postgres_host, postgres_port, postgres_user, postgres_db, postgres_password]):
    print("Ошибка: Не все переменные окружения для подключения к PostgreSQL определены.")
    sys.exit(1)


inserts_filename = "generated_inserts.sql"
with open(inserts_filename, "w") as file:
    file.write("BEGIN;\n")
    file.write(features_inserts[:-2] + ";\n")
    file.write(tags_inserts[:-2] + ";\n")
    file.write("COMMIT;\n")


os.system(f"PGPASSWORD={postgres_password} psql -h {postgres_host} -p {postgres_port} -U {postgres_user} -d {postgres_db} -a -f {inserts_filename}")


os.remove(inserts_filename)

