import pandas as pd
import numpy as np
import pandas_datareader.data as web
import datetime
import matplotlib.pyplot as plt
from scipy.optimize import minimize


"""
S: python markowitz optimization

https://kevinvecmanis.io/finance/optimization/2019/04/02/Algorithmic-Portfolio-Optimization.html

https://github.com/chaitjo/markowitz-portfolio-optimization/

https://realpython.com/fast-flexible-pandas/
https://www.tradewithscience.com/practical-portfolio-optimization-in-python-1-3-markowitz/

https://www.youtube.com/watch?v=57qAxRV577c
Stock Market Analysis with Pandas Python Programming | Python # 6

https://www.youtube.com/watch?v=f2BCmQBCwDs
Stock Market Analysis & Markowitz Efficient Frontier on Python | Python # 11

https://www.youtube.com/watch?v=7kNwJYGghoE
Análisis stock mercado y optimización de cartera Markowitz | Aplicación de optimización convexa #9

https://www.youtube.com/watch?v=BkvyJ0HJAQM
Cree su propio Markowitz Portfolio Solver | Derivaciones | Análisis del mercado de valores

https://quantivity.wordpress.com/2011/02/21/why-log-returns/

https://www.learnpythonwithrune.org/master-markowitz-portfolio-optimization-efficient-frontier-in-python-using-pandas/
"""

if __name__ == '__main__':
    start = datetime.datetime(2021, 1, 1)
    end = datetime.datetime(2021, 1, 21)

    google = web.DataReader('CEVA', 'yahoo', start, end)
    apple = web.DataReader('GOOGL', 'yahoo', start, end)
    nvs = web.DataReader('TSLA', 'yahoo', start, end)
    tesla = web.DataReader('ZOM', 'yahoo', start, end)

    stocks = pd.concat([google['Close'], apple['Close'], nvs['Close'], tesla['Close']], axis=1)

    stocks.columns = ['CEVA', 'GOOGL', 'TSLA', 'ZOM']

    returns = stocks / stocks.shift(1)
    logReturns = np.log(returns)

    noOfPortfolios = 10000
    weights = np.zeros((noOfPortfolios, 4))
    meanLogReturns = logReturns.mean()
    Sigma = logReturns.cov()

    expectedReturn = np.zeros(noOfPortfolios)
    expectedVolatility = np.zeros(noOfPortfolios)
    sharpeRatio = np.zeros(noOfPortfolios)

    for k in range(noOfPortfolios):
        # generate random weights
        w = np.array(np.random.random(4))
        w = w / w.sum()
        weights[k, :] = w

        # expected log return
        expectedReturn[k] = np.sum(meanLogReturns * w)

        # expected volatility
        expectedVolatility[k] = np.sqrt(np.dot(w.T, np.dot(Sigma, w)))

        # sharpe ratio
        sharpeRatio[k] = expectedReturn[k] / expectedVolatility[k]

    maxIndex = sharpeRatio.argmax()
    maxExpectedReturn = expectedReturn.max()
    maxVolatility = expectedVolatility.max()

    def negativeSR(w):
        w = np.array(w)

        R = np.sum(meanLogReturns * w)

        V = np.sqrt(np.dot(w.T, np.dot(Sigma, w)))

        return -R/V

    def checkSumToOne(w):
        return np.sum(w) - 1

    w0 = np.asarray([.25, .25, .25, .25])
    bounds = ((0, 1), (0, 1), (0, 1), (0, 1))
    constraints = ({'type': 'eq', 'fun': checkSumToOne})
    w_opt = minimize(negativeSR, w0, method='SLSQP', bounds=bounds, constraints=constraints)

    print(f'MAX {weights[maxIndex]}')
    print(f'OPT {w_opt.x}')

    returns = np.linspace(0, maxExpectedReturn, 50)
    volatility_opt = []

    def minimizeMyVolatility(w):
        w = np.array(w)
        V = np.sqrt(np.dot(w.T, np.dot(Sigma, w)))
        return V

    def getReturn(w):
        w = np.array(w)

        R = np.sum(meanLogReturns * w)

        return R


    for return__ in returns:
        constraints = ({'type': 'eq', 'fun': checkSumToOne},
                       {'type': 'eq', 'fun': lambda w: getReturn(w) - return__})
        opt = minimize(minimizeMyVolatility, w0, bounds=bounds, constraints=constraints)
        volatility_opt.append(opt['fun'])


    plt.figure(figsize=(12, 8))
    plt.scatter(expectedVolatility, expectedReturn, c=sharpeRatio)
    plt.xlabel('expected volatility')
    plt.ylabel('expected log return')
    plt.colorbar(label='SR')
    plt.scatter(expectedVolatility[maxIndex], expectedReturn[maxIndex], c='red')
    plt.plot(volatility_opt, returns, '--')
    plt.show()
