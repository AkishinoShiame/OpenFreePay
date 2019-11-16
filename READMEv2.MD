{OpenFreePay}  
  
# Concept
  
* 第三方點數兌換戳和平台
  * 定額點數上架
  * 動態戳和點數上架
* 平台會員
  * 手機號碼建會員
    * 簡訊認證（不爽做）
    * 信箱認證（不爽做）
  
# Idea  
  
* 基本頁面:  
  * 網頁版RWD
    
> 使用手機號登入
> 點數兌換刊登頁面（包含定額與不定額）
>> 交易雙方溝通區
> 查詢歷史刊登（含交易紀錄）
  
# Rough MAP
  
* AMAZON RDS (MySQL)
* AMAZON ECR (Docker) 
  * {demonickrace/OpenFreePay}
* AMAZON VPC (NetWork)
  
# Hack 'n' Roll 
![Structure](snapshop.png) 
