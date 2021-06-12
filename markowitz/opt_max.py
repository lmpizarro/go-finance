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

    noOfPortfolios = 100000
    meanLogReturns = logReturns.mean()
    Sigma = logReturns.cov()

    weights = np.zeros((noOfPortfolios, len(tickets)))
    expectedReturn = np.zeros(noOfPortfolios)
    expectedVolatility = np.zeros(noOfPortfolios)
    sharpeRatio = np.zeros(noOfPortfolios)

    print("Sigma")
    print(Sigma)
    print()

    print("log returns")
    print(logReturns)
    print()

    print("mean log returns")
    print(meanLogReturns)
    print()


    for k in range(noOfPortfolios):
        # generate random weights
        w = np.array(np.random.random(len(tickets)))
        w = w / w.sum()
        weights[k, :] = w

        # expected log return
        expectedReturn[k] = np.sum(meanLogReturns * w)

        # expected volatility
        expectedVolatility[k] = np.sqrt(np.dot(w.T, np.dot(Sigma, w)))

        # sharpe ratio
        sharpeRatio[k] = expectedReturn[k] / expectedVolatility[k]

    maxIndex = sharpeRatio.argmax()


    print(f'Return      {expectedReturn[maxIndex]}')
    print(f'Volatility  {expectedVolatility[maxIndex]}')
    print(f'SharpeRatio {sharpeRatio[maxIndex]}')



    print(f'Weights {weights[maxIndex]}')

    print()

    for k,v in zip(tickets, weights[maxIndex]):
        print(f'{k:>10}  {100*v:6.2f}')


    def negativeSR(w):
        w = np.array(w)
        R = np.sum(meanLogReturns * w)
        V = np.sqrt(np.dot(w.T, np.dot(Sigma, w)))
        return -(R-0.001)/V

    def checkSumToOne(w):
        return np.sum(w) - 1

    w0 = np.asarray([.25]*len(tickets))
    bounds = ((0, 1),)*len(tickets)
    constraints = ({'type': 'eq', 'fun': checkSumToOne})
    w_opt = minimize(negativeSR, w0, method='SLSQP', 
                     bounds=bounds, constraints=constraints,
                     options={'disp':True, 'ftol': .0000001})

    print()
    for k,v in zip(tickets, w_opt.x):
        print(f'{k:>10}  {100*v:6.2f}')


    print(negativeSR(w0))

