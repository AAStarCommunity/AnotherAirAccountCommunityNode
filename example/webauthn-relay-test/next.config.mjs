module.exports = {
    async rewrites() {
      return [
        {
          source: '/api/:path*',
          destination: 'https://api.yourservice.com/:path*' // 代理到后端服务
        }
      ];
    }    
  };