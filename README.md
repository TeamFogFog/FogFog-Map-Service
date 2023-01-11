# FogFog-Map-Service
FogFog 내부용 Map Service

### LatAndLong
- Using Go
- Naver Map Api Geocoding 
- Address 가 담긴 xlsx 파일을 가져와 geocoding 하여 위 / 경도를 담은 xlsx 파일을 생성해줌

#### 실행 방법
- .env에 API Key 설정
- open file 설정

``` go
  f, err := excelize.OpenFile("name.xlsx")
  
  ...
  
  rows := f.GetRows("sheet name")
```
- write file name 설정
``` go
newFile.SetCellValue("sheet name", "A"+cellIdx, address)

er := newFile.SaveAs("name.xlsx")
```
