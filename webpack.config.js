const path = require('path')
const fs = require('fs')
const webpack = require('webpack')
const dotenv = require('dotenv')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const CopyWebpackPlugin = require('copy-webpack-plugin')

dotenv.config()

module.exports = {
	entry: {
		chat: path.resolve(__dirname, 'src/index.js'),
	},
	output: {
		path: path.resolve(__dirname, 'dist'),
		filename: '[name].js',
		library: {
			name: '[name]',
			type: 'umd',
		},
	},
	module: {
		rules: [
			{
				test: /\.(js)$/,
				include: path.resolve(__dirname, 'src'),
				// exclude: /node_modules/,
				use: {
					loader: 'babel-loader',
					options: {
						presets: ['@babel/preset-env'],
						plugins: ['@babel/transform-runtime'],
					},
				},
			},
			{
				test: /\.css$/i,
				include: path.resolve(__dirname, 'src'),
				use: [MiniCssExtractPlugin.loader, 'css-loader'],
			},
		],
	},
	resolve: {
		extensions: ['*', '.js'],
	},
	plugins: [
		new MiniCssExtractPlugin({
			filename: '[name].css',
		}),
		new HtmlWebpackPlugin({
			template: path.resolve(__dirname, 'src/index.html'),
			filename: 'index.html',
			inject: 'head',
			scriptLoading: 'blocking',
		}),
		new CopyWebpackPlugin({
			patterns: [
				{
					from: path.resolve(__dirname, 'static'),
					to: 'static',
				},
			],
		}),
		new webpack.DefinePlugin({
			'process.env': JSON.stringify(process.env),
		}),
	],
	devServer: {
		static: {
			directory: path.resolve(__dirname, 'public'),
		},
		port: 9000,
		allowedHosts: ['lo.cal'],
		https: {
			cert: fs.readFileSync('/Users/habin/dev/pki/lo.cal/cert.pem'),
			key: fs.readFileSync('/Users/habin/dev/pki/lo.cal/key.pem'),
		},
	},
}
