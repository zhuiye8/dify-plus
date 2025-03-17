from flask import Flask, request, jsonify
import sqlite3
import os
import json

app = Flask(__name__)

# 禁用JSON中的ASCII转义，确保中文直接显示而不是转义
app.config['JSON_AS_ASCII'] = False

# 数据库设置
DATABASE_PATH = "resume_database.db"

def init_database():
    """初始化数据库，创建多个关联表并添加初始数据"""
    # 检查数据库文件是否存在
    db_exists = os.path.exists(DATABASE_PATH)
    
    conn = sqlite3.connect(DATABASE_PATH)
    conn.text_factory = str  # 确保文本以字符串形式返回
    cursor = conn.cursor()
    
    # 创建基本信息表 - 添加了部门字段
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS persons (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        title TEXT,
        department TEXT,
        phone TEXT,
        email TEXT,
        location TEXT,
        linkedin TEXT,
        website TEXT,
        summary TEXT
    )
    ''')
    
    # 创建工作经历表
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS work_experiences (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        person_id INTEGER NOT NULL,
        job_title TEXT,
        company TEXT,
        location TEXT,
        dates TEXT,
        responsibilities TEXT,
        FOREIGN KEY (person_id) REFERENCES persons (id) ON DELETE CASCADE
    )
    ''')
    
    # 创建教育经历表
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS education (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        person_id INTEGER NOT NULL,
        degree TEXT,
        institution TEXT,
        location TEXT,
        graduation_date TEXT,
        achievements TEXT,
        FOREIGN KEY (person_id) REFERENCES persons (id) ON DELETE CASCADE
    )
    ''')
    
    # 创建技能表
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS skills (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        person_id INTEGER NOT NULL,
        skills TEXT,
        FOREIGN KEY (person_id) REFERENCES persons (id) ON DELETE CASCADE
    )
    ''')
    
    # 创建语言表
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS languages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        person_id INTEGER NOT NULL,
        languages TEXT,
        FOREIGN KEY (person_id) REFERENCES persons (id) ON DELETE CASCADE
    )
    ''')
    
    # 创建证书表
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS certifications (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        person_id INTEGER NOT NULL,
        certifications TEXT,
        FOREIGN KEY (person_id) REFERENCES persons (id) ON DELETE CASCADE
    )
    ''')
    
    # 创建项目表
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS projects (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        person_id INTEGER NOT NULL,
        projects TEXT,
        FOREIGN KEY (person_id) REFERENCES persons (id) ON DELETE CASCADE
    )
    ''')
    
    # 创建志愿者经历表
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS volunteer (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        person_id INTEGER NOT NULL,
        volunteer TEXT,
        FOREIGN KEY (person_id) REFERENCES persons (id) ON DELETE CASCADE
    )
    ''')
    
    # 如果数据库是新创建的，添加初始数据
    if not db_exists:
        # 添加示例数据的格式部分（不包含实际数据）
        cursor.execute('''
        INSERT INTO persons (name, title, department, phone, email, location, linkedin, website, summary)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        ''', (
            '姓名示例',
            '职位示例',
            '部门示例',
            '电话示例',
            'email示例',
            '地点示例',
            'linkedin示例',
            'website示例',
            '简介示例'
        ))
    
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

# 辅助函数：获取完整简历信息
def get_resume_data(person_id, cursor):
    """获取单个人员的完整简历信息"""
    # 获取基本信息
    cursor.execute('SELECT * FROM persons WHERE id = ?', (person_id,))
    person_row = cursor.fetchone()
    
    if not person_row:
        return None
    
    # 构建完整简历数据
    resume = {}
    for key in person_row.keys():
        resume[key] = person_row[key]
    
    # 获取工作经历
    cursor.execute('SELECT * FROM work_experiences WHERE person_id = ? ORDER BY id', (person_id,))
    work_experiences = []
    for row in cursor.fetchall():
        work_exp = {}
        for key in row.keys():
            if key != 'person_id':  # 排除外键
                work_exp[key] = row[key]
        work_experiences.append(work_exp)
    resume['work_experiences'] = work_experiences
    
    # 获取教育经历
    cursor.execute('SELECT * FROM education WHERE person_id = ? ORDER BY id', (person_id,))
    education = []
    for row in cursor.fetchall():
        edu = {}
        for key in row.keys():
            if key != 'person_id':
                edu[key] = row[key]
        education.append(edu)
    resume['education'] = education
    
    # 获取技能
    cursor.execute('SELECT skills FROM skills WHERE person_id = ?', (person_id,))
    skills_row = cursor.fetchone()
    resume['skills'] = skills_row['skills'] if skills_row else None
    
    # 获取语言
    cursor.execute('SELECT languages FROM languages WHERE person_id = ?', (person_id,))
    languages_row = cursor.fetchone()
    resume['languages'] = languages_row['languages'] if languages_row else None
    
    # 获取证书
    cursor.execute('SELECT certifications FROM certifications WHERE person_id = ?', (person_id,))
    certifications_row = cursor.fetchone()
    resume['certifications'] = certifications_row['certifications'] if certifications_row else None
    
    # 获取项目
    cursor.execute('SELECT projects FROM projects WHERE person_id = ?', (person_id,))
    projects_row = cursor.fetchone()
    resume['projects'] = projects_row['projects'] if projects_row else None
    
    # 获取志愿者经历
    cursor.execute('SELECT volunteer FROM volunteer WHERE person_id = ?', (person_id,))
    volunteer_row = cursor.fetchone()
    resume['volunteer'] = volunteer_row['volunteer'] if volunteer_row else None
    
    return resume

# API路由
@app.route('/api/persons', methods=['GET'])
def get_persons():
    """获取所有人员信息，可选择是否包含完整信息"""
    include_full = request.args.get('full', 'false').lower() == 'true'
    name = request.args.get('name', '')
    department = request.args.get('department', '')
    
    conn = sqlite3.connect(DATABASE_PATH)
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()
    
    query = 'SELECT * FROM persons WHERE 1=1'
    params = []
    
    if name:
        query += ' AND name LIKE ?'
        params.append(f'%{name}%')
    
    if department:
        query += ' AND department LIKE ?'
        params.append(f'%{department}%')
    
    cursor.execute(query, params)
    
    result = []
    for row in cursor.fetchall():
        if include_full:
            # 获取完整信息
            resume = get_resume_data(row['id'], cursor)
            if resume:
                result.append(resume)
        else:
            # 仅获取基本信息
            person = {}
            for key in row.keys():
                person[key] = row[key]
            result.append(person)
    
    conn.close()
    return jsonify(result)

@app.route('/api/persons/<int:person_id>', methods=['GET'])
def get_person_complete(person_id):
    """获取单个人员的完整简历信息（包括所有关联表数据）"""
    conn = sqlite3.connect(DATABASE_PATH)
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()
    
    resume = get_resume_data(person_id, cursor)
    
    conn.close()
    
    if not resume:
        return jsonify({'error': '人员不存在'}), 404
    
    return jsonify(resume)

@app.route('/api/persons', methods=['POST'])
def add_person():
    """添加新人员及其完整简历信息"""
    data = request.json
    
    if not data or 'name' not in data:
        return jsonify({'error': '姓名是必填项'}), 400
    
    conn = sqlite3.connect(DATABASE_PATH)
    cursor = conn.cursor()
    
    try:
        # 开始事务
        conn.execute('BEGIN')
        
        # 添加基本信息（新增部门字段）
        cursor.execute('''
        INSERT INTO persons (name, title, department, phone, email, location, linkedin, website, summary)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
        ''', (
            data.get('name'),
            data.get('title'),
            data.get('department'),
            data.get('phone'),
            data.get('email'),
            data.get('location'),
            data.get('linkedin'),
            data.get('website'),
            data.get('summary')
        ))
        
        person_id = cursor.lastrowid
        
        # 添加工作经历
        if 'job_title_1' in data and 'company_1' in data:
            cursor.execute('''
            INSERT INTO work_experiences (person_id, job_title, company, location, dates, responsibilities)
            VALUES (?, ?, ?, ?, ?, ?)
            ''', (
                person_id,
                data.get('job_title_1'),
                data.get('company_1'),
                data.get('location_1'),
                data.get('dates_1'),
                data.get('responsibilities_1')
            ))
        
        if 'job_title_2' in data and 'company_2' in data:
            cursor.execute('''
            INSERT INTO work_experiences (person_id, job_title, company, location, dates, responsibilities)
            VALUES (?, ?, ?, ?, ?, ?)
            ''', (
                person_id,
                data.get('job_title_2'),
                data.get('company_2'),
                data.get('location_2'),
                data.get('dates_2'),
                data.get('responsibilities_2')
            ))
        
        if 'job_title_3' in data and 'company_3' in data:
            cursor.execute('''
            INSERT INTO work_experiences (person_id, job_title, company, location, dates, responsibilities)
            VALUES (?, ?, ?, ?, ?, ?)
            ''', (
                person_id,
                data.get('job_title_3'),
                data.get('company_3'),
                data.get('location_3'),
                data.get('dates_3'),
                data.get('responsibilities_3')
            ))
        
        # 添加教育经历
        if 'degree_1' in data and 'institution_1' in data:
            cursor.execute('''
            INSERT INTO education (person_id, degree, institution, location, graduation_date, achievements)
            VALUES (?, ?, ?, ?, ?, ?)
            ''', (
                person_id,
                data.get('degree_1'),
                data.get('institution_1'),
                data.get('edu_location_1'),
                data.get('graduation_date_1'),
                data.get('achievements_1')
            ))
        
        if 'degree_2' in data and 'institution_2' in data:
            cursor.execute('''
            INSERT INTO education (person_id, degree, institution, location, graduation_date, achievements)
            VALUES (?, ?, ?, ?, ?, ?)
            ''', (
                person_id,
                data.get('degree_2'),
                data.get('institution_2'),
                data.get('edu_location_2'),
                data.get('graduation_date_2'),
                data.get('achievements_2')
            ))
        
        # 添加技能
        if 'skills' in data:
            cursor.execute('''
            INSERT INTO skills (person_id, skills)
            VALUES (?, ?)
            ''', (person_id, data.get('skills')))
        
        # 添加语言
        if 'languages' in data:
            cursor.execute('''
            INSERT INTO languages (person_id, languages)
            VALUES (?, ?)
            ''', (person_id, data.get('languages')))
        
        # 添加证书
        if 'certifications' in data:
            cursor.execute('''
            INSERT INTO certifications (person_id, certifications)
            VALUES (?, ?)
            ''', (person_id, data.get('certifications')))
        
        # 添加项目
        if 'projects' in data:
            cursor.execute('''
            INSERT INTO projects (person_id, projects)
            VALUES (?, ?)
            ''', (person_id, data.get('projects')))
        
        # 添加志愿者经历
        if 'volunteer' in data:
            cursor.execute('''
            INSERT INTO volunteer (person_id, volunteer)
            VALUES (?, ?)
            ''', (person_id, data.get('volunteer')))
        
        # 提交事务
        conn.commit()
        
    except Exception as e:
        # 回滚事务
        conn.rollback()
        conn.close()
        return jsonify({'error': f'添加失败: {str(e)}'}), 500
    
    conn.close()
    return jsonify({'id': person_id, 'message': '简历添加成功'}), 201

@app.route('/api/persons/<int:person_id>', methods=['PUT'])
def update_person(person_id):
    """更新人员简历信息"""
    data = request.json
    
    if not data:
        return jsonify({'error': '没有提供更新数据'}), 400
    
    conn = sqlite3.connect(DATABASE_PATH)
    cursor = conn.cursor()
    
    try:
        # 开始事务
        conn.execute('BEGIN')
        
        # 检查人员是否存在
        cursor.execute('SELECT id FROM persons WHERE id = ?', (person_id,))
        if not cursor.fetchone():
            conn.close()
            return jsonify({'error': '人员不存在'}), 404
        
        # 更新基本信息
        fields = []
        values = []
        
        for key in ['name', 'title', 'department', 'phone', 'email', 'location', 'linkedin', 'website', 'summary']:
            if key in data:
                fields.append(f"{key} = ?")
                values.append(data.get(key))
        
        if fields:
            values.append(person_id)
            query = f"UPDATE persons SET {', '.join(fields)} WHERE id = ?"
            cursor.execute(query, values)
        
        # 更新工作经历（这里简化处理，实际应用中可能需要更复杂的逻辑）
        # 先删除现有记录，再添加新记录
        if any(key.startswith('job_title_') for key in data):
            cursor.execute('DELETE FROM work_experiences WHERE person_id = ?', (person_id,))
            
            if 'job_title_1' in data and 'company_1' in data:
                cursor.execute('''
                INSERT INTO work_experiences (person_id, job_title, company, location, dates, responsibilities)
                VALUES (?, ?, ?, ?, ?, ?)
                ''', (
                    person_id,
                    data.get('job_title_1'),
                    data.get('company_1'),
                    data.get('location_1'),
                    data.get('dates_1'),
                    data.get('responsibilities_1')
                ))
            
            if 'job_title_2' in data and 'company_2' in data:
                cursor.execute('''
                INSERT INTO work_experiences (person_id, job_title, company, location, dates, responsibilities)
                VALUES (?, ?, ?, ?, ?, ?)
                ''', (
                    person_id,
                    data.get('job_title_2'),
                    data.get('company_2'),
                    data.get('location_2'),
                    data.get('dates_2'),
                    data.get('responsibilities_2')
                ))
            
            if 'job_title_3' in data and 'company_3' in data:
                cursor.execute('''
                INSERT INTO work_experiences (person_id, job_title, company, location, dates, responsibilities)
                VALUES (?, ?, ?, ?, ?, ?)
                ''', (
                    person_id,
                    data.get('job_title_3'),
                    data.get('company_3'),
                    data.get('location_3'),
                    data.get('dates_3'),
                    data.get('responsibilities_3')
                ))
        
        # 更新教育经历
        if any(key.startswith('degree_') for key in data):
            cursor.execute('DELETE FROM education WHERE person_id = ?', (person_id,))
            
            if 'degree_1' in data and 'institution_1' in data:
                cursor.execute('''
                INSERT INTO education (person_id, degree, institution, location, graduation_date, achievements)
                VALUES (?, ?, ?, ?, ?, ?)
                ''', (
                    person_id,
                    data.get('degree_1'),
                    data.get('institution_1'),
                    data.get('edu_location_1'),
                    data.get('graduation_date_1'),
                    data.get('achievements_1')
                ))
            
            if 'degree_2' in data and 'institution_2' in data:
                cursor.execute('''
                INSERT INTO education (person_id, degree, institution, location, graduation_date, achievements)
                VALUES (?, ?, ?, ?, ?, ?)
                ''', (
                    person_id,
                    data.get('degree_2'),
                    data.get('institution_2'),
                    data.get('edu_location_2'),
                    data.get('graduation_date_2'),
                    data.get('achievements_2')
                ))
        
        # 更新技能
        if 'skills' in data:
            cursor.execute('UPDATE skills SET skills = ? WHERE person_id = ?', 
                          (data.get('skills'), person_id))
            if cursor.rowcount == 0:  # 如果没有更新任何行（即该技能不存在）
                cursor.execute('INSERT INTO skills (person_id, skills) VALUES (?, ?)',
                              (person_id, data.get('skills')))
        
        # 更新语言
        if 'languages' in data:
            cursor.execute('UPDATE languages SET languages = ? WHERE person_id = ?', 
                          (data.get('languages'), person_id))
            if cursor.rowcount == 0:  # 如果没有更新任何行
                cursor.execute('INSERT INTO languages (person_id, languages) VALUES (?, ?)',
                              (person_id, data.get('languages')))
        
        # 更新证书
        if 'certifications' in data:
            cursor.execute('UPDATE certifications SET certifications = ? WHERE person_id = ?', 
                          (data.get('certifications'), person_id))
            if cursor.rowcount == 0:
                cursor.execute('INSERT INTO certifications (person_id, certifications) VALUES (?, ?)',
                              (person_id, data.get('certifications')))
        
        # 更新项目
        if 'projects' in data:
            cursor.execute('UPDATE projects SET projects = ? WHERE person_id = ?', 
                          (data.get('projects'), person_id))
            if cursor.rowcount == 0:
                cursor.execute('INSERT INTO projects (person_id, projects) VALUES (?, ?)',
                              (person_id, data.get('projects')))
        
        # 更新志愿者经历
        if 'volunteer' in data:
            cursor.execute('UPDATE volunteer SET volunteer = ? WHERE person_id = ?', 
                          (data.get('volunteer'), person_id))
            if cursor.rowcount == 0:
                cursor.execute('INSERT INTO volunteer (person_id, volunteer) VALUES (?, ?)',
                              (person_id, data.get('volunteer')))
        
        # 提交事务
        conn.commit()
        
    except Exception as e:
        # 回滚事务
        conn.rollback()
        conn.close()
        return jsonify({'error': f'更新失败: {str(e)}'}), 500
    
    conn.close()
    return jsonify({'message': '简历更新成功'})

@app.route('/api/persons/<int:person_id>', methods=['DELETE'])
def delete_person(person_id):
    """删除人员及其所有相关信息"""
    conn = sqlite3.connect(DATABASE_PATH)
    cursor = conn.cursor()
    
    # 检查人员是否存在
    cursor.execute('SELECT id FROM persons WHERE id = ?', (person_id,))
    if not cursor.fetchone():
        conn.close()
        return jsonify({'error': '人员不存在'}), 404
    
    try:
        # 开始事务
        conn.execute('BEGIN')
        
        # 删除所有关联表中的数据
        cursor.execute('DELETE FROM work_experiences WHERE person_id = ?', (person_id,))
        cursor.execute('DELETE FROM education WHERE person_id = ?', (person_id,))
        cursor.execute('DELETE FROM skills WHERE person_id = ?', (person_id,))
        cursor.execute('DELETE FROM languages WHERE person_id = ?', (person_id,))
        cursor.execute('DELETE FROM certifications WHERE person_id = ?', (person_id,))
        cursor.execute('DELETE FROM projects WHERE person_id = ?', (person_id,))
        cursor.execute('DELETE FROM volunteer WHERE person_id = ?', (person_id,))
        
        # 删除人员基本信息
        cursor.execute('DELETE FROM persons WHERE id = ?', (person_id,))
        
        # 提交事务
        conn.commit()
        
    except Exception as e:
        # 回滚事务
        conn.rollback()
        conn.close()
        return jsonify({'error': f'删除失败: {str(e)}'}), 500
    
    conn.close()
    return jsonify({'message': '简历删除成功'})

@app.route('/api/persons/search', methods=['GET'])
def search_persons():
    """搜索人员简历"""
    query = request.args.get('q', '')
    name = request.args.get('name', '')
    department = request.args.get('department', '')
    include_full = request.args.get('full', 'false').lower() == 'true'
    
    conn = sqlite3.connect(DATABASE_PATH)
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()
    
    # 构建搜索条件
    search_conditions = []
    params = []
    
    if query:
        search_conditions.extend([
            "p.name LIKE ?",
            "p.title LIKE ?",
            "p.department LIKE ?",
            "p.summary LIKE ?",
            "w.job_title LIKE ?",
            "w.company LIKE ?",
            "e.degree LIKE ?",
            "e.institution LIKE ?",
            "s.skills LIKE ?"
        ])
        query_param = f'%{query}%'
        params.extend([query_param] * 9)
    
    if name:
        search_conditions.append("p.name LIKE ?")
        params.append(f'%{name}%')
    
    if department:
        search_conditions.append("p.department LIKE ?")
        params.append(f'%{department}%')
    
    # 如果没有搜索条件，默认返回所有
    if not search_conditions:
        cursor.execute('SELECT id, name, title, department, email FROM persons')
    else:
        sql = '''
        SELECT DISTINCT p.id, p.name, p.title, p.department, p.email
        FROM persons p
        LEFT JOIN work_experiences w ON p.id = w.person_id
        LEFT JOIN education e ON p.id = e.person_id
        LEFT JOIN skills s ON p.id = s.person_id
        WHERE ''' + ' OR '.join(search_conditions)
        cursor.execute(sql, params)
    
    # 处理结果
    result = []
    for row in cursor.fetchall():
        if include_full:
            # 获取完整信息
            resume = get_resume_data(row['id'], cursor)
            if resume:
                result.append(resume)
        else:
            # 仅获取基本信息
            person = {}
            for key in row.keys():
                person[key] = row[key]
            result.append(person)
    
    conn.close()
    return jsonify(result)

@app.route('/api/persons/by-departments', methods=['GET'])
def get_persons_by_department():
    """获取按部门分组的人员信息"""
    include_full = request.args.get('full', 'false').lower() == 'true'
    
    conn = sqlite3.connect(DATABASE_PATH)
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()
    
    # 获取所有部门
    cursor.execute('SELECT DISTINCT department FROM persons WHERE department IS NOT NULL')
    departments = [row['department'] for row in cursor.fetchall()]
    
    result = {}
    
    # 对每个部门，获取其人员信息
    for dept in departments:
        cursor.execute('SELECT * FROM persons WHERE department = ?', (dept,))
        
        dept_persons = []
        for row in cursor.fetchall():
            if include_full:
                # 获取完整信息
                resume = get_resume_data(row['id'], cursor)
                if resume:
                    dept_persons.append(resume)
            else:
                # 仅获取基本信息
                person = {}
                for key in row.keys():
                    person[key] = row[key]
                dept_persons.append(person)
        
        result[dept] = dept_persons
    
    conn.close()
    return jsonify(result)

# 修复中文编码的响应头设置
@app.after_request
def after_request(response):
    response.headers.add('Content-Type', 'application/json; charset=utf-8')
    return response

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=9000)