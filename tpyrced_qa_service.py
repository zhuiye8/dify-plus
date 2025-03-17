from flask import Flask, request, jsonify
import sqlite3
import os
import json

app = Flask(__name__)

# 禁用JSON中的ASCII转义，确保中文直接显示
app.config['JSON_AS_ASCII'] = False

# 数据库设置
DATABASE_PATH = "qa_database.db"

def init_database():
    """初始化数据库，创建问答表"""
    conn = sqlite3.connect(DATABASE_PATH)
    conn.text_factory = str
    cursor = conn.cursor()
    
    # 创建问答表
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS qas (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        question TEXT NOT NULL,
        answer TEXT NOT NULL
    )
    ''')
    
    conn.commit()
    conn.close()

# 初始化数据库
init_database()

# 自定义JSON编码器处理中文
class ChineseJSONEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, bytes):
            return obj.decode('utf-8')
        return json.JSONEncoder.default(self, obj)

app.json_encoder = ChineseJSONEncoder

# API路由

@app.route('/api/qas', methods=['GET'])
def get_qas():
    """获取所有问答记录"""
    conn = sqlite3.connect(DATABASE_PATH)
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()
    
    question = request.args.get('question', '')
    
    if question:
        cursor.execute('SELECT * FROM qas WHERE question LIKE ?', (f'%{question}%',))
    else:
        cursor.execute('SELECT * FROM qas')
    
    result = [dict(row) for row in cursor.fetchall()]
    
    conn.close()
    return jsonify(result)

@app.route('/api/qas/<int:qa_id>', methods=['GET'])
def get_qa(qa_id):
    """获取单个问答记录"""
    conn = sqlite3.connect(DATABASE_PATH)
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()
    
    cursor.execute('SELECT * FROM qas WHERE id = ?', (qa_id,))
    qa = cursor.fetchone()
    
    conn.close()
    
    if not qa:
        return jsonify({'error': '问答记录不存在'}), 404
    
    return jsonify(dict(qa))

@app.route('/api/qas', methods=['POST'])
def add_qa():
    """添加单个或批量问答记录"""
    data = request.json
    
    conn = sqlite3.connect(DATABASE_PATH)
    cursor = conn.cursor()
    
    try:
        conn.execute('BEGIN')
        
        if isinstance(data, list):  # 批量添加
            inserted_ids = []
            for item in data:
                if 'question' not in item or 'answer' not in item:
                    conn.rollback()
                    conn.close()
                    return jsonify({'error': '每个记录必须包含问题和答案'}), 400
                
                cursor.execute('''
                INSERT INTO qas (question, answer)
                VALUES (?, ?)
                ''', (item['question'], item['answer']))
                inserted_ids.append(cursor.lastrowid)
            
            conn.commit()
            conn.close()
            return jsonify({'ids': inserted_ids, 'message': '批量问答添加成功'}), 201
            
        else:  # 单个添加
            if 'question' not in data or 'answer' not in data:
                conn.rollback()
                conn.close()
                return jsonify({'error': '问题和答案是必填项'}), 400
            
            cursor.execute('''
            INSERT INTO qas (question, answer)
            VALUES (?, ?)
            ''', (data['question'], data['answer']))
            
            qa_id = cursor.lastrowid
            conn.commit()
            conn.close()
            return jsonify({'id': qa_id, 'message': '问答添加成功'}), 201
            
    except Exception as e:
        conn.rollback()
        conn.close()
        return jsonify({'error': f'添加失败: {str(e)}'}), 500

@app.route('/api/qas/<int:qa_id>', methods=['PUT'])
def update_qa(qa_id):
    """更新单个问答记录"""
    data = request.json
    
    if not data:
        return jsonify({'error': '没有提供更新数据'}), 400
    
    conn = sqlite3.connect(DATABASE_PATH)
    cursor = conn.cursor()
    
    try:
        conn.execute('BEGIN')
        
        cursor.execute('SELECT id FROM qas WHERE id = ?', (qa_id,))
        if not cursor.fetchone():
            conn.close()
            return jsonify({'error': '问答记录不存在'}), 404
        
        fields = []
        values = []
        
        for key in ['question', 'answer']:
            if key in data:
                fields.append(f"{key} = ?")
                values.append(data[key])
        
        if fields:
            values.append(qa_id)
            query = f"UPDATE qas SET {', '.join(fields)} WHERE id = ?"
            cursor.execute(query, values)
        
        conn.commit()
        conn.close()
        return jsonify({'message': '问答更新成功'})
        
    except Exception as e:
        conn.rollback()
        conn.close()
        return jsonify({'error': f'更新失败: {str(e)}'}), 500

@app.route('/api/qas/batch-update', methods=['PUT'])
def batch_update_qa():
    """批量更新问答记录"""
    data = request.json
    
    if not isinstance(data, list):
        return jsonify({'error': '批量更新需要提供记录数组'}), 400
    
    conn = sqlite3.connect(DATABASE_PATH)
    cursor = conn.cursor()
    
    try:
        conn.execute('BEGIN')
        
        updated_ids = []
        for item in data:
            if 'id' not in item:
                conn.rollback()
                conn.close()
                return jsonify({'error': '每条记录必须包含ID'}), 400
                
            cursor.execute('SELECT id FROM qas WHERE id = ?', (item['id'],))
            if not cursor.fetchone():
                conn.rollback()
                conn.close()
                return jsonify({'error': f'ID {item["id"]} 不存在'}), 404
            
            fields = []
            values = []
            for key in ['question', 'answer']:
                if key in item:
                    fields.append(f"{key} = ?")
                    values.append(item[key])
            
            if fields:
                values.append(item['id'])
                query = f"UPDATE qas SET {', '.join(fields)} WHERE id = ?"
                cursor.execute(query, values)
                updated_ids.append(item['id'])
        
        conn.commit()
        conn.close()
        return jsonify({'updated_ids': updated_ids, 'message': '批量更新成功'})
        
    except Exception as e:
        conn.rollback()
        conn.close()
        return jsonify({'error': f'批量更新失败: {str(e)}'}), 500

@app.route('/api/qas/<int:qa_id>', methods=['DELETE'])
def delete_qa(qa_id):
    """删除单个问答记录"""
    conn = sqlite3.connect(DATABASE_PATH)
    cursor = conn.cursor()
    
    try:
        conn.execute('BEGIN')
        
        cursor.execute('SELECT id FROM qas WHERE id = ?', (qa_id,))
        if not cursor.fetchone():
            conn.close()
            return jsonify({'error': '问答记录不存在'}), 404
        
        cursor.execute('DELETE FROM qas WHERE id = ?', (qa_id,))
        
        conn.commit()
        conn.close()
        return jsonify({'message': '问答删除成功'})
        
    except Exception as e:
        conn.rollback()
        conn.close()
        return jsonify({'error': f'删除失败: {str(e)}'}), 500

@app.route('/api/qas/batch-delete', methods=['DELETE'])
def batch_delete_qa():
    """批量删除问答记录"""
    data = request.json
    
    if not data or 'ids' not in data or not isinstance(data['ids'], list):
        return jsonify({'error': '需要提供要删除的ID数组'}), 400
    
    conn = sqlite3.connect(DATABASE_PATH)
    cursor = conn.cursor()
    
    try:
        conn.execute('BEGIN')
        
        deleted_ids = []
        for qa_id in data['ids']:
            cursor.execute('SELECT id FROM qas WHERE id = ?', (qa_id,))
            if cursor.fetchone():
                cursor.execute('DELETE FROM qas WHERE id = ?', (qa_id,))
                deleted_ids.append(qa_id)
        
        conn.commit()
        conn.close()
        return jsonify({'deleted_ids': deleted_ids, 'message': '批量删除成功'})
        
    except Exception as e:
        conn.rollback()
        conn.close()
        return jsonify({'error': f'批量删除失败: {str(e)}'}), 500

@app.after_request
def after_request(response):
    response.headers.add('Content-Type', 'application/json; charset=utf-8')
    return response

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=9001)