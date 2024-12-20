import pandas as pd
import matplotlib.pyplot as plt

# Загрузка данных из файла
file_path = 'stats1.log'
data = pd.read_csv(file_path)

# Преобразуем столбец Timestamp в формат datetime
data['Timestamp'] = pd.to_datetime(data['Timestamp'])

# Устанавливаем Timestamp как индекс
data.set_index('Timestamp', inplace=True)

# График общего количества запросов во времени
plt.figure(figsize=(12, 6))
plt.plot(data.index, data['TotalRequests'], label='Total Requests')
plt.title('Total Requests Over Time')
plt.xlabel('Timestamp')
plt.ylabel('Total Requests')
plt.legend()
plt.grid()
plt.show()

# График количества отклонённых запросов во времени
plt.figure(figsize=(12, 6))
plt.plot(data.index, data['RejectedRequests'], label='Rejected Requests', color='red')
plt.title('Rejected Requests Over Time')
plt.xlabel('Timestamp')
plt.ylabel('Rejected Requests')
plt.legend()
plt.grid()
plt.show()

# График вероятности отклонения запросов во времени
plt.figure(figsize=(12, 6))
plt.plot(data.index, data['ProbabilityOfRejection'], label='Probability of Rejection', color='orange')
plt.title('Probability of Rejection Over Time')
plt.xlabel('Timestamp')
plt.ylabel('Probability of Rejection')
plt.legend()
plt.grid()
plt.show()

# График среднего времени ожидания в буфере и среднего времени обработки во времени
plt.figure(figsize=(12, 6))
plt.plot(data.index, data['AverageBufferTime'], label='Average Buffer Time', color='green')
plt.plot(data.index, data['AverageProcessingTime'], label='Average Processing Time', color='blue')
plt.title('Average Buffer Time and Processing Time Over Time')
plt.xlabel('Timestamp')
plt.ylabel('Time (seconds)')
plt.legend()
plt.grid()
plt.show()

# График загрузки специалистов во времени
specialist_columns = [col for col in data.columns if 'Specialist' in col]
plt.figure(figsize=(12, 6))
for col in specialist_columns:
    plt.plot(data.index, data[col], label=col)
plt.title('Specialist Load Over Time')
plt.xlabel('Timestamp')
plt.ylabel('Load')
plt.legend()
plt.grid()
plt.show()