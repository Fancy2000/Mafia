import sqlite3
from flask import jsonify, url_for

DB_PATH = './todo.db'   # Update this path accordingly
NOTSTARTED = 'Not Started'


#DONE--------------
# Формирование структуры ответа на основе информации, представленной в базе данных
def make_public_users(row):
    new_todo = {}
    for field in row.keys():
        print(field)
        # Предоставление URL объекта вместо его id - хорошая практика
        if field == 'user_name':
            new_todo['uri'] = url_for('get_user', user_name=row['user_name'], _external=True)
        else:
            new_todo[field] = row[field]
    print("AFter for")
    print(new_todo)
    return new_todo


#DONE--------------
# Получить все элементы в таблице
def get_all_users():
    try:
        conn = sqlite3.connect(DB_PATH)
        # Обеспечивает работу с названиями колонок в таблице
        conn.row_factory = sqlite3.Row
        c = conn.cursor()
        c.execute('select * from users')
        # Получаем список строк в перечислимом формате
        conn.commit()
        rows = c.fetchall()
        # С помощью функции map применяем функцию make_public_todo ко всем элементам rows
        result = jsonify({'users': list(map(make_public_users, rows))})
        return result
    except Exception as e:
        print('Error: ', e)
        return None


#DONE------------
# Получить отдельного пользователя
def get_user(user_name):
    try:
        print("in get user")
        conn = sqlite3.connect(DB_PATH)
        # Обеспечивает работу с названиями колонок в таблице
        conn.row_factory = sqlite3.Row
        c = conn.cursor()
        c.execute("select * from users where user_name=?;", [user_name])
        conn.commit()
        r = c.fetchone()
        print(r)
        return jsonify(make_public_users(r))
    except Exception as e:
        print('Error: ', e)
    return None


#DONE-----------------------
# Добавить элемент в таблицу
def add_to_list(user_name, avatar, sex, email):
    try:
        conn = sqlite3.connect(DB_PATH)
        c = conn.cursor()
        c.execute('insert into users(user_name, avatar, sex, email) values(?,?, ?, ?)', (user_name, avatar, sex, email))
        conn.commit()
        result = get_user(user_name)
        print(result)

        return result
    except Exception as e:
        print('Error: ', e)
        return None


#DONE---------------
# Обновить элемент с user_name в таблице
def update_users(user_name, avatar, sex, email):
    try:
        conn = sqlite3.connect(DB_PATH)
        c = conn.cursor()
        c.execute('update users set avatar=?, sex=?, email=? where user_name=?', (avatar, sex, email, user_name))
        conn.commit()
        result = get_user(user_name)
        return result
    except Exception as e:
        print('Error: ', e)
        return None


#DONE--------------------------
# Удалить элемент из таблицы по имени
def remove_user(user_name):
    try:
        conn = sqlite3.connect(DB_PATH)
        c = conn.cursor()
        c.execute('DELETE FROM users WHERE user_name=?', [user_name])
        conn.commit()
        return jsonify( { 'result': True } )
    except Exception as e:
        print('Error: ', e)
        return None