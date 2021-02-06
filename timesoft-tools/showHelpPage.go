package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func showHelpPage() fyne.CanvasObject {
	usage := `
使用帮助：

Excel相关操作：
	Trim Space工具：可以将输入文件内所有表格的所有单元格前后的空白去除
	身份证号码扩展：可以将输入文件内所有的15位数字扩展为18位身份证号
	案卷目录格式检测：可以检查输入文件文件是否符合案卷目录规范

文件夹重命名：
	可将输入的项目文件夹根据 档案号+姓名 在所选xlsx文件的第一个Sheet中匹配数据
	并批量将文件夹重命名为对应身份证号。

	*注意！
	1. 所选xlsx文件的第一个Sheet中应有两行标题
	2. 工号/档案号 应在第三列
	3. 姓名 应在第四列
	4. 身份证号 应在第七列
	5. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*

文件重命名：
	可将 输入的项目文件夹 下全部文件及全部子文件重命名为 文件上级文件夹名+原文件名 + (1 ；
	若文件夹名为“目录”，则将此文件夹下文件重命名为 “目录”上级文件夹名 + 0( + 以罗马数字3开始的序号

	*注意！
	1. 程序可递归查找全部的文件夹及子文件夹
	2. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*

封皮&档案袋重命名：
	可将 输入的装有 imageXXXX.jpg 的文件夹 下全部 imageXXXX.jpg 文件
	按顺序依次重命名为所选xlsx文档第一个Sheet中的第三列（身份证号）。

	*注意！
	1. 案卷目录.xlsx 的第一个Sheet中应有两行标题
	2. 程序仅查找输入的文件夹，不会递归查找子文件夹
	3. 身份证号 应在第三列
	4. 应保证身份证号列的身份证号数量与读取到的jpg文件数量完全一致
	5. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*

封皮&档案袋归位：
	可将输入的装有 身份证号.jpg 的文件夹下全部 身份证号.jpg 文件
	移动到输入的项目文件夹中对应的身份证号文件夹的“目录”内。

	*注意！
	1. 对于包含jpg文件的文件夹，程序不会递归查找子文件夹
	2. 对于项目文件夹，程序将递归查找子文件夹
	3. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*

移动目录：
	可将选中的文件夹下全部名为“目录”的子文件夹内的文件转移至其上层文件夹。

	*注意！
	1. 程序会递归查找名称不为“目录”的全部子文件夹
	2. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*


Author: 3RM1N3@时源科技
E-mail: wangyu7439@hotmail.com`
	return container.NewBorder(
		nil,
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewLabel("版权所有 © 2021 上海时源信息科技有限公司。保留所有权利。\nCopyright © 2021 Timesoft Corporation. All Rights Reserved."),
			layout.NewSpacer(),
		),
		nil,
		nil,
		container.NewVScroll(widget.NewLabel(usage)),
	)
}
