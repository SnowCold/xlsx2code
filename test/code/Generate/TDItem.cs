//Auto Generate Don't Edit it
using UnityEngine;
using System;
using System.IO;
using System.Collections;
using System.Collections.Generic;
namespace LR
{
    public partial class TDItem
    {
        
       
        private string m_Id;   
        private string m_ItemName;   
        private string m_ItemDes;   
        private EInt m_Rare = 0;   
        private EInt m_ItemLevel = 0;   
        private EInt m_Type = 0;   
        private EInt m_StackNum = 0;   
        private EInt m_NumMax = 0;   
        private EInt m_UsedType = 0;   
        private string m_FlowID;  
        
        private Dictionary<string, TDUniversally.FieldData> m_DataCacheNoGenerate = new Dictionary<string, TDUniversally.FieldData>();
      
        /// <summary>
        /// 道具ID
        /// </summary>
        public  string  id {get { return m_Id; } }
       
        /// <summary>
        /// 道具名字
        /// </summary>
        public  string  itemName {get { return m_ItemName; } }
       
        /// <summary>
        /// 道具描述
        /// </summary>
        public  string  itemDes {get { return m_ItemDes; } }
       
        /// <summary>
        /// 道具品质
        /// </summary>
        public  int  rare {get { return m_Rare; } }
       
        /// <summary>
        /// 道具等级
        /// </summary>
        public  int  itemLevel {get { return m_ItemLevel; } }
       
        /// <summary>
        /// 道具类型
        /// </summary>
        public  int  type {get { return m_Type; } }
       
        /// <summary>
        /// 叠加数量
        /// </summary>
        public  int  stackNum {get { return m_StackNum; } }
       
        /// <summary>
        /// 可拥有数量上限
        /// </summary>
        public  int  numMax {get { return m_NumMax; } }
       
        /// <summary>
        /// 使用类型
        /// </summary>
        public  int  usedType {get { return m_UsedType; } }
       
        /// <summary>
        /// 流程ID
        /// </summary>
        public  string  flowID {get { return m_FlowID; } }
       

        public void ReadRow(DataStreamReader dataR, int[] filedIndex)
        {
          var schemeNames = dataR.GetSchemeName();
          int col = 0;
          while(true)
          {
            col = dataR.MoreFieldOnRow();
            if (col == -1)
            {
              break;
            }
            switch (filedIndex[col])
            { 
            
                case 0:
                    m_Id = dataR.ReadString();
                    break;
                case 1:
                    m_ItemName = dataR.ReadString();
                    break;
                case 2:
                    m_ItemDes = dataR.ReadString();
                    break;
                case 3:
                    m_Rare = dataR.ReadInt();
                    break;
                case 4:
                    m_ItemLevel = dataR.ReadInt();
                    break;
                case 5:
                    m_Type = dataR.ReadInt();
                    break;
                case 6:
                    m_StackNum = dataR.ReadInt();
                    break;
                case 7:
                    m_NumMax = dataR.ReadInt();
                    break;
                case 8:
                    m_UsedType = dataR.ReadInt();
                    break;
                case 9:
                    m_FlowID = dataR.ReadString();
                    break;
                default:
                    TableHelper.CacheNewField(dataR, schemeNames[col], m_DataCacheNoGenerate);
                    break;
            }
          }

        }
        
        public DataStreamReader.FieldType GetFieldTypeInNew(string fieldName)
        {
            if (m_DataCacheNoGenerate.ContainsKey(fieldName))
            {
                return m_DataCacheNoGenerate[fieldName].fieldType;
            }
            return DataStreamReader.FieldType.Unkown;
        }
        
        public static Dictionary<string, int> GetFieldHeadIndex()
        {
          Dictionary<string, int> ret = new Dictionary<string, int>(10);
          
          ret.Add("Id", 0);
          ret.Add("ItemName", 1);
          ret.Add("ItemDes", 2);
          ret.Add("Rare", 3);
          ret.Add("ItemLevel", 4);
          ret.Add("Type", 5);
          ret.Add("StackNum", 6);
          ret.Add("NumMax", 7);
          ret.Add("UsedType", 8);
          ret.Add("FlowID", 9);
          return ret;
        }
    } 
}//namespace LR