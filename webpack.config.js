const path = require('path')
const fs = require('fs')
const HtmlWebpackPlugin = require('html-webpack-plugin')

module.exports = {
	entry: path.resolve(__dirname, 'index.js'),
	output: {
		path: path.resolve(__dirname, 'dist'),
		filename: 'bundle.js',
	},
	module: {
		rules: [
			{
				test: /\.(js)$/,
				exclude: /node_modules/,
				use: {
					loader: 'babel-loader',
					options: {
						presets: ['@babel/preset-env', '@babel/preset-react'],
					},
				},
			},
		],
	},
	resolve: {
		extensions: ['*', '.js'],
	},
	plugins: [
		new HtmlWebpackPlugin({
			template: path.resolve(__dirname, 'index.html'),
		}),
	],
	devServer: {
		static: {
			directory: path.resolve(__dirname, 'public'),
		},
		port: 8000,
		allowedHosts: ['lo.cal'],
		https: {
			cert: fs.readFileSync('/Users/habin/dev/pki/lo.cal/cert.pem'),
			key: fs.readFileSync('/Users/habin/dev/pki/lo.cal/key.pem'),
		},
	},
}
