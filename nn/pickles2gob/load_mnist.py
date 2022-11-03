import pickle

# 重点是rb和r的区别，rb是打开2进制文件，文本文件用r
f = open('mnist.pkl','rb')

u = pickle._Unpickler( f )
u.encoding = 'latin1'

training_data, validation_data, test_data = u.load()
print(training_data)

f.close()

