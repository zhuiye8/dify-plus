import requests
import random
from faker import Faker
import json

# 初始化 Faker，支持中文数据
fake = Faker('zh_CN')

# API 基础 URL
BASE_URL = "http://0.0.0.0:9000/api/persons"

# 部门列表
DEPARTMENTS = ["技术部", "产品部", "数据部", "市场部", "人力资源部"]

# 职位列表
TITLES = ["软件工程师", "产品经理", "数据分析师", "市场专员", "人力资源专员", "高级开发工程师"]

# 技能列表
SKILLS = ["Python", "Java", "SQL", "React", "Excel", "数据分析", "项目管理", "市场营销"]

# 语言列表
LANGUAGES = ["中文", "英文", "日文", "法文"]

# 证书列表
CERTIFICATIONS = ["PMP", "CFA", "AWS Certified", "CPA", "RHCE"]

# 随机生成一条简历数据
def generate_fake_resume():
    name = fake.name()  # 随机中文姓名
    title = random.choice(TITLES)
    department = random.choice(DEPARTMENTS)
    
    resume = {
        "name": name,
        "title": title,
        "department": department,
        "phone": fake.phone_number(),
        "email": fake.email(),
        "location": fake.city(),
        "linkedin": f"linkedin.com/in/{name.lower().replace(' ', '')}",
        "website": f"{name.lower().replace(' ', '')}.com",
        "summary": fake.text(max_nb_chars=100),
        
        # 随机生成 1-3 条工作经历
        "job_title_1": random.choice(TITLES),
        "company_1": fake.company(),
        "location_1": fake.city(),
        "dates_1": f"{random.randint(2015, 2020)}-{random.randint(2021, 2025)}",
        "responsibilities_1": fake.text(max_nb_chars=150),
    }
    
    # 随机决定是否添加第二条工作经历
    if random.random() > 0.5:
        resume.update({
            "job_title_2": random.choice(TITLES),
            "company_2": fake.company(),
            "location_2": fake.city(),
            "dates_2": f"{random.randint(2010, 2015)}-{random.randint(2016, 2020)}",
            "responsibilities_2": fake.text(max_nb_chars=150),
        })
    
    # 随机决定是否添加第三条工作经历
    if random.random() > 0.7:
        resume.update({
            "job_title_3": random.choice(TITLES),
            "company_3": fake.company(),
            "location_3": fake.city(),
            "dates_3": f"{random.randint(2005, 2010)}-{random.randint(2011, 2015)}",
            "responsibilities_3": fake.text(max_nb_chars=150),
        })
    
    # 教育经历
    resume.update({
        "degree_1": random.choice(["学士", "硕士", "博士"]) + "学位",
        "institution_1": fake.company() + "大学",
        "edu_location_1": fake.city(),
        "graduation_date_1": str(random.randint(2000, 2023)),
        "achievements_1": fake.text(max_nb_chars=50),
    })
    
    # 技能、语言、证书、项目、志愿者经历
    resume["skills"] = ", ".join(random.sample(SKILLS, random.randint(2, 5)))
    resume["languages"] = ", ".join(random.sample(LANGUAGES, random.randint(1, 3)))
    resume["certifications"] = ", ".join(random.sample(CERTIFICATIONS, random.randint(1, 3)))
    resume["projects"] = ", ".join([fake.bs() for _ in range(random.randint(1, 3))])
    resume["volunteer"] = fake.text(max_nb_chars=100)
    
    return resume

# 添加数据到 API
def add_fake_data(count=5):
    for _ in range(count):
        resume = generate_fake_resume()
        try:
            response = requests.post(BASE_URL, json=resume)
            if response.status_code == 201:
                print(f"成功添加: {resume['name']} (ID: {response.json()['id']})")
            else:
                print(f"添加失败: {resume['name']} - {response.status_code} - {response.text}")
        except Exception as e:
            print(f"请求出错: {str(e)}")

if __name__ == "__main__":
    # 生成并添加 10 条假数据
    add_fake_data(count=10)