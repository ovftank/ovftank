import re
from datetime import date
birthday = date(2004, 5, 19)
today = date.today()
age = today.year - birthday.year
if today.month < birthday.month or (today.month == birthday.month and today.day < birthday.day):
    age -= 1
with open('README.md', 'r') as file:
    content = file.read()
new_content = re.sub(r'AGE: \d+', f'AGE: {age}', content)
with open('README.md', 'w') as file:
    file.write(new_content)