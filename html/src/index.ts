/**
 * NATS.js 浏览器客户端演示
 * 
 * 这个模块实现了一个基于浏览器的NATS客户端应用，支持：
 * - 手动连接/断开NATS服务器
 * - JetStream消息消费
 * - 普通主题订阅和发布
 * - 请求/响应模式
 * 
 * @author Zhenjiang Zhang
 * @version 1.0.0
 */

'use strict';
import * as nats from "nats.ws"
//import * as nats from "https://cdn.skypack.dev/nats.ws";

/** NATS连接类型 */
type NatsConnection = nats.NatsConnection;
/** 订阅类型 */
type Subscription = nats.Subscription;
/** 消息类型 */
type Msg = nats.Msg;


/** NATS错误接口，扩展标准Error */
interface NatsError extends Error {
	/** 错误代码 */
	code?: string;
}

/** 全局NATS连接实例 */
let natsConnection: NatsConnection | null = null;
/** 全局字符串编解码器实例 */
let stringCodec: nats.Codec<string> | null = null;

/**
 * 连接到NATS服务器
 * 建立WebSocket连接，初始化JetStream管理器，设置订阅和事件处理
 * @returns Promise<void> 连接完成后的Promise
 */
async function connectToNats(): Promise<void> {
	try {
		updateStatus('连接中...');
		updateConnectButtons(false, true);

		// 获取用户选择的服务器地址和token
		const serverUrl = getSelectedServerUrl();
		const token = getAuthToken();

		if (!serverUrl) {
			updateStatus('请选择或输入有效的服务器地址');
			updateConnectButtons(false, false);
			return;
		}

		const nc: NatsConnection = await nats.connect({
			servers: serverUrl,
			token: token
		});
		console.log('连接到 NATS 服务器:', serverUrl);

		natsConnection = nc;
		stringCodec = nats.StringCodec();
		const js = nc.jetstream({ domain: "hub" });

		// 初始化 JetStream 管理器
		await initJetStreamManager(js, stringCodec);

		// 设置订阅
		setupSubscriptions(nc, stringCodec);

		// 设置资源清理
		setupCleanup(nc);

		updateStatus('已连接');
		updateConnectButtons(true, false);
		enablePublishButton(true);

	} catch (err) {
		console.error('连接 NATS 服务器失败:', err);
		updateStatus('连接失败');
		updateConnectButtons(false, false);
	}
}

/**
 * 断开与NATS服务器的连接
 * 关闭连接，清理全局变量，更新UI状态
 * @returns Promise<void> 断开连接完成后的Promise
 */
async function disconnectFromNats(): Promise<void> {
	if (natsConnection) {
		try {
			await natsConnection.close();
			natsConnection = null;
			stringCodec = null;
			updateStatus('已断开连接');
			updateConnectButtons(false, false);
			enablePublishButton(false);
			console.log('已断开 NATS 连接');
		} catch (err) {
			console.error('断开连接时出错:', err);
		}
	}
}

/**
 * 初始化JetStream管理器
 * 创建消费者，设置消息消费回调
 * @param js JetStream客户端实例
 * @param sc 字符串编解码器
 * @returns Promise<void> 初始化完成后的Promise
 */
async function initJetStreamManager(js: nats.JetStreamClient, sc: nats.Codec<string>): Promise<void> {
	try {
		const jsm = await js.jetstreamManager(false);

		// 获取消费者列表
		const consumers: nats.ConsumerInfo[] = await jsm.consumers.list("EVENTS").next();
		console.log('已经获取 JetStream 消费者列表:');
		consumers.forEach((consumer) => {
			console.log("item:", consumer.name);
		});
		// await jsm.consumers.delete("EVENTS", "Hello");
		jsm.consumers.add("EVENTS", {
			name: "Hello",
			description: "Hello consumer",
			"ack_policy": nats.AckPolicy.Explicit,
			"durable_name": "Hello",
		}).then((consumerInfo: nats.ConsumerInfo) => {
			console.log('已创建 JetStream 消费者:', consumerInfo.name);
		}).catch((err: Error) => {
			console.error('创建 JetStream 消费者时出错:', err);
		});
		// console.log('已创建 JetStream 消费者');
		// // 更新消费者
		// const consumerInfo = await jsm.consumers.update("EVENTS", "Hello", {
		// 	description: "Hello consumer",
		// });
		// console.log('已更新 JetStream 消费者:', consumerInfo.name);

		// 获取并消费消息
		const consumer: nats.Consumer = await js.consumers.get("EVENTS", "Hello");
		const sub: nats.ConsumerMessages = await consumer.consume({
			callback: (msg: nats.JsMsg) => {
				try {
					// Try to parse the message as JSON
					const msgData = JSON.parse(sc.decode(msg.data));
					console.log('Successfully parsed JSON message:', msgData);

				} catch (jsonError) {
					console.log('Message is not valid JSON, treating as plain text:', sc.decode(msg.data));
				} finally {
					msg.ack();
				};

				addMessageToUI("SUBJECT:"+msg.subject+" MSG:"+sc.decode(msg.data) + " SEQ:" + msg.seq);
			}
		});
		console.log('已经打开 JetStream 消消费者');

	} catch (err) {
		console.error('JetStream 操作错误:', err);
	}
}


/**
 * 设置普通主题订阅
 * 订阅demo.topic和demo.cb主题，处理接收到的消息
 * @param nc NATS连接实例
 * @param sc 字符串编解码器
 */
function setupSubscriptions(nc: NatsConnection, sc: nats.Codec<string>): void {
	// 订阅消息
	nc.subscribe("demo.topic", {
		callback: (err: NatsError | null, msg: Msg) => {
			if (err) {
				console.error('订阅错误:', err);
				return;
			}
			const msgText: string = sc.decode(msg.data);
			addMessageToUI(msgText);
		}
	});

	// 回调服务
	nc.subscribe("demo.cb", {
		callback: (err: NatsError | null, msg: Msg) => {
			if (err) {
				console.error('订阅错误:', err);
				return;
			}
			const msgText: string = sc.decode(msg.data);
			addMessageToUI(msgText);
			msg.respond(sc.encode('Hello from browser at ' + new Date().toLocaleTimeString()));
		}
	});
}

/**
 * 更新状态显示
 * @param status 要显示的状态文本
 */
function updateStatus(status: string): void {
	const statusElement = document.getElementById('status');
	if (statusElement) {
		statusElement.textContent = status;
	}
}

/**
 * 更新连接和断开连接按钮的状态
 * @param disconnectEnabled 是否启用断开连接按钮
 * @param connectDisabled 是否禁用连接按钮
 */
function updateConnectButtons(disconnectEnabled: boolean, connectDisabled: boolean): void {
	const connectButton = document.getElementById('connect') as HTMLButtonElement;
	const disconnectButton = document.getElementById('disconnect') as HTMLButtonElement;

	if (connectButton) {
		connectButton.disabled = connectDisabled;
	}
	if (disconnectButton) {
		disconnectButton.disabled = !disconnectEnabled;
	}
}

/**
 * 启用或禁用发布消息按钮
 * @param enabled 是否启用发布按钮
 */
function enablePublishButton(enabled: boolean): void {
	const publishButton = document.getElementById('publish') as HTMLButtonElement;
	const requestButton = document.getElementById('request-btn') as HTMLButtonElement;
	if (publishButton) {
		publishButton.disabled = !enabled;
	}
	if (requestButton) {
		requestButton.disabled = !enabled;
	}
}

/**
 * 设置UI事件处理器
 * 为连接、断开连接和发布消息按钮绑定点击事件
 */
function setupUIHandlers(): void {
	// 服务器地址选择变化处理
	const serverSelect = document.getElementById('nats-server-select');
	if (serverSelect) {
		serverSelect.onchange = handleServerSelectChange;
	}

	// 连接按钮
	const connectButton = document.getElementById('connect');
	if (connectButton) {
		connectButton.onclick = (): void => {
			connectToNats();
		};
	}

	// 断开连接按钮
	const disconnectButton = document.getElementById('disconnect');
	if (disconnectButton) {
		disconnectButton.onclick = (): void => {
			disconnectFromNats();
		};
	}

	// 发布消息按钮
	const publishButton = document.getElementById('publish');
	if (publishButton) {
		publishButton.onclick = (): void => {
							if (natsConnection && stringCodec) {
					const msg: string = 'Hello from browser at ' + new Date().toLocaleTimeString();
					natsConnection.publish("demo.topic", stringCodec.encode(msg));
				}
			};
		}

		// 请求按钮
		const requestButton = document.getElementById('request-btn');
		if (requestButton) {
			requestButton.onclick = async (): Promise<void> => {
				const subject = (document.getElementById('request-subject') as HTMLInputElement).value;
				const payload = (document.getElementById('request-payload') as HTMLInputElement).value;
				if (natsConnection && stringCodec) {
					try {
						// Display that we're sending the request
						addMessageToUI(`Sending request to ${subject} with payload: ${payload}`);
						
						// Send the request and wait for a response
						const response = await natsConnection.request(
							subject,
							stringCodec.encode(payload),
							{ timeout: 2000 } // 2 second timeout
						);
						
						// Display the response
						const responseText = stringCodec.decode(response.data);
						addMessageToUI(`Received response: ${responseText}`);
					} catch (err) {
						console.error('Request failed:', err);
						addMessageToUI(`Request failed: ${err instanceof Error ? err.message : String(err)}`);
					}
				} else {
					addMessageToUI('Not connected to NATS server');
				}
			};
		}
	}


/**
 * 将消息添加到UI显示列表
 * @param msgText 要显示的消息文本
 */
function addMessageToUI(msgText: string): void {
	const li: HTMLLIElement = document.createElement('li');
	li.textContent = '收到消息: ' + msgText;
	const messagesElement = document.getElementById('messages');
	if (messagesElement) {
		messagesElement.appendChild(li);
		// 限制消息数量，避免页面卡顿
		if (messagesElement.children.length > 100) {
			messagesElement.removeChild(messagesElement.firstChild!);
		}
	}
	console.log('收到消息:', msgText);
}

/**
 * 设置资源清理
 * 在页面卸载时自动关闭NATS连接
 * @param nc NATS连接实例
 */
function setupCleanup(nc: NatsConnection): void {
	window.addEventListener('beforeunload', () => {
		nc.close();
	});
}

/**
 * 获取用户选择的服务器地址
 * @returns string | null 返回选择的服务器地址，如果无效则返回null
 */
function getSelectedServerUrl(): string | null {
	const selectElement = document.getElementById('nats-server-select') as HTMLSelectElement;
	const customInputElement = document.getElementById('custom-server-input') as HTMLInputElement;
	
	if (!selectElement) {
		return null;
	}
	
	const selectedValue = selectElement.value;
	
	if (selectedValue === 'custom') {
		const customUrl = customInputElement?.value?.trim();
		if (!customUrl) {
			return null;
		}
		// 简单验证WebSocket URL格式
		if (!customUrl.startsWith('ws://') && !customUrl.startsWith('wss://')) {
			return null;
		}
		return customUrl;
	}
	
	return selectedValue;
}

/**
 * 获取认证token
 * @returns string 返回用户输入的token，如果为空则返回默认值
 */
function getAuthToken(): string {
	const tokenElement = document.getElementById('token-input') as HTMLInputElement;
	return tokenElement?.value?.trim() || '';
}

/**
 * 处理服务器地址选择变化
 */
function handleServerSelectChange(): void {
	const selectElement = document.getElementById('nats-server-select') as HTMLSelectElement;
	const customInputElement = document.getElementById('custom-server-input') as HTMLInputElement;
	
	if (!selectElement || !customInputElement) {
		return;
	}
	
	if (selectElement.value === 'custom') {
		customInputElement.style.display = 'inline-block';
		customInputElement.focus();
	} else {
		customInputElement.style.display = 'none';
	}
}

// 应用初始化
document.addEventListener('DOMContentLoaded', () => {
	setupUIHandlers();
	updateStatus('未连接');
	updateConnectButtons(false, false);
	enablePublishButton(false);
});
