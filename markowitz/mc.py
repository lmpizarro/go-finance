from numpy.core.fromnumeric import ptp
import pandas as pd
import numpy as np
import pandas_datareader.data as web
import datetime
from scipy.optimize import minimize


if __name__ == '__main__':
    start = datetime.datetime(2021, 5, 10)
    end = datetime.datetime(2021, 6, 10)
    tickets = ['CEVA', 'GOOGL', 'TSLA', 'ZOM']
    tickets = ['FB', 'GOOGL', 'AAPL', 'AMZN',] 
    # 'TSLA', 'DIS', 'NVS', 'NVDA', 'TSM', 'KO', 'TXN', 'AMD']


    print("get data")
    columns = []

    for ticket in tickets:
        data = web.DataReader(ticket, 'yahoo', start, end)
        columns.append(data['Close'])

    stocks = pd.concat(columns, axis=1)
    stocks.columns = tickets 

    print('end get data')

    returns = stocks / stocks.shift(1)
    logReturns = np.log(returns)

    meanLogReturns = logReturns.mean()
    Sigma = logReturns.cov()

    print("Sigma")
    print(Sigma)
    print()

    print("log returns")
    print(logReturns)
    print()

    print("mean log returns")
    print(meanLogReturns)
    print()


    noOfPortfolios = 40000
    weights = np.zeros((noOfPortfolios, len(tickets)))
    expectedReturn = np.zeros(noOfPortfolios)
    expectedVolatility = np.zeros(noOfPortfolios)
    sharpeRatio = np.zeros(noOfPortfolios)

    def negativeSR(w):
        w = np.array(w)
        R = np.sum(meanLogReturns * w)
        V = np.sqrt(np.dot(w.T, np.dot(Sigma, w)))
        return -(R-0.001)/V

    E0 = 0
    
    K = 100.0
    Ts = [.01 - i/3000.0 for i in range(30)]
    Ts = 0.01*np.exp(-np.asarray(range(30))/5)
    print(Ts)
    w0 = [0]* len(tickets)
    ET = 0
    for T in Ts:
        for k in range(noOfPortfolios):
            # generate random weights
            w =  np.random.rand(len(tickets))
            
            w = w / w.sum()
            
            E = negativeSR(w)
            DE = E - E0 
            if DE < 0:
                w0 = w
                E0 = E
            elif DE > 0:
                pacc = np.exp(-DE/T)
                r = np.random.random()
                if pacc > r:
            
                    w0 = w
                    E0 = E

        print(ET, E0, np.abs(ET - E0),  w0)
        ET = E0
    
    

